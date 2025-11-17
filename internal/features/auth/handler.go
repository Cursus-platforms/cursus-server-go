package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Cursus-platforms/cursus-server-go/internal/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

type registerRequest struct {
	Fullname string
	Email    string
	Password string
}

// @Summary      Register a new user
// @Description  Creates a new user account with fullname, email, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        registerBody  body  registerRequest  true  "User registration information"
// @Success      201  {object}  model.User  "Account created successfully (Returns user info WITHOUT password)"
// @Failure      400  {object}  map[string]string "Invalid request body"
// @Failure      409  {object}  map[string]string "Email already exists"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /auth/register [post]
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.svc.Register(r.Context(), req.Fullname, req.Email, req.Password)
	if err != nil {
		if err == ErrUserExisted {
			response.RespondWithError(w, http.StatusConflict, err.Error())
		} else {
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, user)
}
