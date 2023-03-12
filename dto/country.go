package dto

type CreateCountryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdateCountryRequest struct {
	Name string `json:"name" form:"name"`
}

type CountryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
