package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/httputil"
	v "github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/validator"
	"github.com/Cursus-platforms/cursus-server-go/internal/lib/response"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

// HandleRequestVerification @Summary Yêu cầu mã xác minh Email
// @Description Nhận thông tin đăng ký, gửi mã OTP và lưu tạm thời dữ liệu user vào Redis.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        requestBody  body  registerRequest  true  "Thông tin đăng ký"
// @Success      202  {object}  map[string]string "OTP đã được gửi thành công"
// @Failure      400  {object}  response.StructuredErrorResponse "Lỗi validation"
// @Failure      409  {object}  map[string]string "Email đã tồn tại"
// @Router       /api/v1/auth/register [post]
func (h *Handler) HandleRequestVerification(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body format", nil)
		return
	}

	if err := v.GlobalValidator.Struct(req); err != nil {
		validationErrors := v.FormatValidationErrors(err)
		response.RespondWithError(w, http.StatusBadRequest, "Input validation failed", validationErrors)
		return
	}

	err := h.svc.RequestVerification(r.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserExisted) {
			response.RespondWithError(w, http.StatusConflict, err.Error(), nil)
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to send verification code", nil)
		return
	}

	response.RespondWithJSON(w, http.StatusAccepted, map[string]string{"message": "OTP sent to email. Please verify."})
}

// HandleVerifyAndRegister @Summary Xác minh OTP và hoàn tất đăng ký
// @Description Xác minh OTP nhận được và lưu dữ liệu người dùng vào Postgres vĩnh viễn.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        verifyBody  body  verifyRequest  true  "Email và mã OTP"
// @Success      201  {object}  map[string]string "Tài khoản được tạo thành công"
// @Failure      401  {object}  map[string]string "OTP không hợp lệ/hết hạn"
// @Failure      500  {object}  map[string]string "Lỗi server nội bộ"
// @Router       /api/v1/auth/verify [post]
func (h *Handler) HandleVerifyAndRegister(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body format", nil)
		return
	}

	if err := v.GlobalValidator.Struct(req); err != nil {
		validationErrors := v.FormatValidationErrors(err)
		response.RespondWithError(w, http.StatusBadRequest, "Input validation failed", validationErrors)
		return
	}

	user, err := h.svc.VerifyAndRegister(r.Context(), req.Email, req.OTP)
	if err != nil {
		if errors.Is(err, ErrInvalidOTP) {
			response.RespondWithError(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, "Registration finalization failed", nil)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, user)
}

// HandleLogin @Summary User login
// @Description  Authenticates user credentials, returns Access Token (JSON body), and sets Refresh Token (HTTP-Only Cookie)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginBody  body  loginRequest  true  "User email and password"
// @Success      200  {object}  loginResponse "Login successful: Access Token in JSON, Refresh Token in Cookie"
// @Failure      401  {object}  map[string]string "Invalid email or password"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/v1/auth/login [post]
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body format", nil)
		return
	}

	accessToken, refreshToken, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.RespondWithError(w, http.StatusUnauthorized, err.Error(), nil)
		} else {
			response.RespondWithError(w, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	httputil.SetRefreshTokenCookie(w, refreshToken, RefreshTokenDuration)

	response.RespondWithJSON(w, http.StatusOK, loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	})
}

// HandleRefresh @Summary      Refresh Access Token
// @Description  Uses the Refresh Token from the HTTP-Only Cookie to issue a new Access Token (JSON) and a new Refresh Token (Cookie)
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  loginResponse "Returns new Access Token (JSON) and sets a new Refresh Token in the Cookie"
// @Failure      401  {object}  map[string]string "Refresh Token is missing, expired, or revoked"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/v1/auth/refresh [post]
func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "Missing refresh token.", nil)
		return
	}

	refreshTokenString := cookie.Value

	newAT, newRT, err := h.svc.Refresh(r.Context(), refreshTokenString)
	if err != nil {
		httputil.SetRefreshTokenCookie(w, "", -time.Hour) // Xóa cookie cũ
		response.RespondWithError(w, http.StatusUnauthorized, "Refresh token is invalid or expired", nil)
		return
	}

	httputil.SetRefreshTokenCookie(w, newRT, RefreshTokenDuration)

	response.RespondWithJSON(w, http.StatusOK, loginResponse{
		AccessToken: newAT,
		TokenType:   "Bearer",
	})
}
