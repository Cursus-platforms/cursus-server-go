package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Cursus-platforms/cursus-server-go/internal/features/role"
	"github.com/Cursus-platforms/cursus-server-go/internal/features/user"
	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

var (
	ErrUserExisted        = errors.New("user aldready exists with this mail")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

var jwtSecretKey = []byte("")

type Service interface {
	Register(ctx context.Context, fullname, email, password string) (model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type service struct {
	userRepo user.Repository
	roleRepo role.Repository
}

func NewService(userRepo user.Repository, roleRepo role.Repository) Service {
	return &service{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *service) Register(ctx context.Context, fullname, email, password string) (model.User, error) {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return model.User{}, ErrUserExisted
	}
	if err != nil && err != sql.ErrNoRows {
		return model.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	hashedPasswordString := string(hashedPassword)

	studentRole, err := s.roleRepo.GetByName(ctx, model.ROLE_STUDENT)
	if err != nil {
		return model.User{}, errors.New("default role not found")
	}

	newUser := model.User{
		Fullname: &fullname,
		Email:    &email,
		Password: &hashedPasswordString,
		RoleID:   &studentRole.ID,
	}

	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return model.User{}, err
	}

	createdUser.Password = nil
	return createdUser, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	foundUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*foundUser.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":  foundUser.ID,
		"role": foundUser.RoleID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
