package course

type createCourseRequest struct {
	Title         string  `json:"title" validate:"required"`
	Description   *string `json:"description,omitempty"`
	Price         float64 `json:"price" validate:"required,gt=0"`
	SubCategoryID string  `json:"sub_category_id" validate:"required,uuid"`
}
