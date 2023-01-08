package countriesdto

type CountryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name" form:"name"`
}
