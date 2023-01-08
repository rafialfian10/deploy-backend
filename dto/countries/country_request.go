package countriesdto

type CreateCountryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdateCountryRequest struct {
	Name string `json:"name" form:"name"`
}
