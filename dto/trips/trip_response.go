package tripsdto

import "project/models"

type TripResponse struct {
	Id             int                    `json:"id"`
	Title          string                 `json:"title"`
	CountryId      int                    `json:"country_id"`
	Country        models.CountryResponse `json:"country"`
	Accomodation   string                 `json:"accomodation"`
	Transportation string                 `json:"transportation"`
	Eat            string                 `json:"eat"`
	Day            int                    `json:"day"`
	Night          int                    `json:"night"`
	DateTrip       string                 `json:"datetrip"`
	Price          int                    `json:"price"`
	Quota          int                    `json:"quota"`
	Description    string                 `json:"description"`
	Image          string                 `json:"image"`
}
