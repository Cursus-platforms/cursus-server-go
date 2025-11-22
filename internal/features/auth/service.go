package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/mailer"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/Cursus-platforms/cursus-server-go/internal/features/role"
	"github.com/Cursus-platforms/cursus-server-go/internal/features/user"
	"github.com/Cursus-platforms/cursus-server-go/internal/lib/utils"
	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

var (
	ErrUserExisted         = errors.New("user already exists with this mail")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

var jwtSecretKey = []byte("")

type Service interface {
	RequestVerification(ctx context.Context, fullname, email, password string) error
	VerifyAndRegister(ctx context.Context, email, otp string) (model.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, refreshTokenString string) (string, string, error)
}

type service struct {
	userRepo user.Repository
	roleRepo role.Repository
	rdb      *redis.Client
	mailer   mailer.Mailer
}

func NewService(userRepo user.Repository, roleRepo role.Repository, rdb *redis.Client, mailer mailer.Mailer) Service {
	return &service{
		userRepo: userRepo,
		roleRepo: roleRepo,
		rdb:      rdb,
		mailer:   mailer,
	}
}

type TempUser struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"` // Mật khẩu đã hash
	OTP      string `json:"otp"`
}

const OtpDuration = 5 * time.Minute // OTP chỉ có giá trị 5 phút

func (s *service) RequestVerification(ctx context.Context, fullname, email, password string) error {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return ErrUserExisted
		}
		return err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	otp := utils.GenerateRandomOTP(6)

	tempUser := TempUser{
		FullName: fullname,
		Email:    email,
		Password: string(hashedPassword),
		OTP:      otp,
	}

	redisKey := fmt.Sprintf("verify:%s", email)
	tempUserData, _ := json.Marshal(tempUser)

	err = s.rdb.Set(ctx, redisKey, tempUserData, OtpDuration).Err()
	if err != nil {
		return fmt.Errorf("lỗi lưu Redis: %w", err)
	}

	log.Printf("DEBUG: Sending OTP %s to %s", otp, email)
	if err := s.mailer.SendOTP(email, fullname, otp); err != nil {
		log.Printf("Error sending email: %v", err)
		// Dọn dẹp key Redis nếu gửi mail thất bại
		s.rdb.Del(ctx, fmt.Sprintf("verify:%s", email))
		return errors.New("lỗi gửi email xác minh")
	}

	return nil
}

var ErrInvalidOTP = errors.New("invalid or expired OTP")

func (s *service) VerifyAndRegister(ctx context.Context, email, otp string) (model.User, error) {
	redisKey := fmt.Sprintf("verify:%s", email)

	// 1. LẤY DỮ LIỆU TỪ REDIS
	val, err := s.rdb.Get(ctx, redisKey).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return model.User{}, ErrInvalidOTP
	}
	if err != nil {
		return model.User{}, fmt.Errorf("lỗi đọc Redis: %w", err)
	}

	var tempUser TempUser
	if err := json.Unmarshal([]byte(val), &tempUser); err != nil {
		return model.User{}, fmt.Errorf("lỗi Unmarshal dữ liệu tạm: %w", err)
	}

	// 2. SO SÁNH OTP
	if tempUser.OTP != otp {
		return model.User{}, ErrInvalidOTP
	}

	// 3. XÁC MINH THÀNH CÔNG: LƯU VÀO POSTGRES
	studentRole, _ := s.roleRepo.GetByName(ctx, model.ROLE_STUDENT)

	newUser := model.User{
		FullName: &tempUser.FullName,
		Email:    &tempUser.Email,
		Password: &tempUser.Password, // Mật khẩu đã hash
		RoleID:   &studentRole.ID,
		IsActive: boolPtr(true),
	}

	registeredUser, err := s.userRepo.Create(ctx, newUser)

	// 4. DỌN DẸP REDIS
	s.rdb.Del(ctx, redisKey)

	return registeredUser, err
}

// Hàm tiện ích đơn giản
func boolPtr(b bool) *bool { return &b }

func (s *service) GenerateTokenPair(userID uuid.UUID, roleName string) (string, string, error) {
	jti := uuid.New().String()

	atClaims := jwt.MapClaims{
		"sub":  userID,
		"role": roleName,
		"exp":  time.Now().Add(AccessTokenDuration).Unix(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	rtClaims := jwt.MapClaims{
		"sub":  userID,
		"role": roleName,
		"jti":  jti, // <-- ID duy nhất để xác minh/thu hồi
		"exp":  time.Now().Add(RefreshTokenDuration).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	redisKey := fmt.Sprintf("rt:%s:%s", jti, userID.String())
	s.rdb.Set(context.Background(), redisKey, "valid", RefreshTokenDuration)

	return accessToken, refreshToken, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, string, error) {
	foundUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*foundUser.Password), []byte(password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	var roleName = model.ROLE_STUDENT
	if foundUser.RoleID != nil {
		roles, err := s.roleRepo.GetByID(ctx, *foundUser.RoleID)
		if err == nil {
			roleName = roles.Name
		}
	}

	return s.GenerateTokenPair(foundUser.ID, roleName)
}

func (s *service) Refresh(ctx context.Context, refreshTokenString string) (string, string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) { return jwtSecretKey, nil })
	if err != nil || !token.Valid {
		return "", "", ErrInvalidRefreshToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	jti, _ := claims["jti"].(string)
	userIDStr, _ := claims["sub"].(string)
	roleName, _ := claims["role"].(string)

	redisKey := fmt.Sprintf("rt:%s:%s", jti, userIDStr)
	status, err := s.rdb.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) || status != "valid" {
		return "", "", ErrInvalidRefreshToken
	}

	s.rdb.Del(ctx, redisKey)

	userID, _ := uuid.Parse(userIDStr)
	return s.GenerateTokenPair(userID, roleName)
}
