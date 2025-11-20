package course

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Cursus-platforms/cursus-server-go/internal/lib/response"
	"github.com/Cursus-platforms/cursus-server-go/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	v "github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/validator"
	validator "github.com/go-playground/validator/v10"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

type createCourseRequest struct {
	Title         string  `json:"title" validate:"required"`
	Description   *string `json:"description,omitempty"`
	Price         float64 `json:"price" validate:"required,gt=0"`
	SubCategoryID string  `json:"sub_category_id" validate:"required,uuid"`
}

func formatValidationErrors(err error) []response.ValidationErrorDetail {
	var details []response.ValidationErrorDetail
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range fieldErrors {
			details = append(details, response.ValidationErrorDetail{
				Field: e.Field(),
				Error: fmt.Sprintf("Field '%s' failed on the '%s' tag", e.Field(), e.Tag()),
			})
		}
	}
	return details
}

// @Summary      Tạo khóa học mới
// @Description  Tạo khóa học mới (yêu cầu Instructor/Admin)
// @Tags         Course
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        requestBody  body  createCourseRequest  true  "Thông tin khóa học mới"
// @Success      201  {object}  model.Course
// @Failure      400  {object}  response.StructuredErrorResponse "Lỗi validation"
// @Router       /api/v1/courses [post]
func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req createCourseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request body format", nil)
		return
	}

	if err := v.GlobalValidator.Struct(req); err != nil {
		validationErrors := formatValidationErrors(err)
		response.RespondWithError(w, http.StatusBadRequest, "Input validation failed", validationErrors)
		return
	}

	userIDStr := r.Context().Value(middleware.UserIDKey).(string)
	userID, _ := uuid.Parse(userIDStr)

	subCatID, _ := uuid.Parse(req.SubCategoryID)

	course := model.Course{
		UserID:        userID,
		Title:         req.Title,
		Price:         req.Price,
		SubCategoryID: &subCatID,
		Description:   req.Description,
	}

	// 4. Gọi Service
	createdCourse, err := h.svc.CreateCourse(r.Context(), course)
	if err != nil {
		if errors.Is(err, ErrPriceRequired) {
			response.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to create course", nil)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, createdCourse)
}

// @Summary      Lấy thông tin chi tiết khóa học
// @Description  Trả về khóa học theo ID (công khai)
// @Tags         Course
// @Produce      json
// @Param        id   path      string  true  "ID của khóa học (UUID)"
// @Success      200  {object}  model.Course
// @Failure      400  {object}  response.StructuredErrorResponse "Invalid ID format"
// @Failure      404  {object}  response.StructuredErrorResponse "Course not found"
// @Router       /api/v1/courses/{id} [get]
func (h *Handler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(idStr)

	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid course ID format", nil)
		return
	}

	course, err := h.svc.GetCourseByID(r.Context(), courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithError(w, http.StatusNotFound, "Course not found", nil)
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve course", nil)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, course)
}
