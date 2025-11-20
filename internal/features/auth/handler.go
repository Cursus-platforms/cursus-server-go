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

// Register handler
type registerRequest struct {
	FullName string `json:"fullname" validate:"required"`
	Email    string `json:"password" validate:"required,min=8,max=50"`
	Password string `json:"email" validate:"required,email"`
}

// HandleRegister @Summary      Register a new user
// @Description  Creates a new user account with fullname, email, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        registerBody  body  registerRequest  true  "User registration information"
// @Success      201  {object}  model.User  "Account created successfully (Returns user info WITHOUT password)"
// @Failure      400  {object}  map[string]string "Invalid request body"
// @Failure      409  {object}  map[string]string "Email already exists"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/v1/auth/register [post]
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := v.GlobalValidator.Struct(req); err != nil {
		// Tạo mảng lỗi chi tiết để trả về theo best practice
		validationErrors := formatValidationErrors(err) // Hàm helper mới

		response.RespondWithError(w, http.StatusBadRequest, "Input validation failed", validationErrors)
		return
	}

	user, err := h.svc.VerifyAndRegister(r.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserExisted) {
			response.RespondWithError(w, http.StatusConflict, err.Error())
		} else {
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, user)
}

// Login handler
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// HandleLogin @Summary      User login
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
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	accessToken, refreshToken, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			response.RespondWithError(w, http.StatusInternalServerError, err.Error())
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
		response.RespondWithError(w, http.StatusUnauthorized, "Missing refresh token.")
		return
	}

	refreshTokenString := cookie.Value

	newAT, newRT, err := h.svc.Refresh(r.Context(), refreshTokenString)
	if err != nil {
		httputil.SetRefreshTokenCookie(w, "", -time.Hour)
		response.RespondWithError(w, http.StatusUnauthorized, "Refresh token is invalid or expired")
		return
	}

	httputil.SetRefreshTokenCookie(w, newRT, RefreshTokenDuration)

	response.RespondWithJSON(w, http.StatusOK, loginResponse{
		AccessToken: newAT,
		TokenType:   "Bearer",
	})
}
