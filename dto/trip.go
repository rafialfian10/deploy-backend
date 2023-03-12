package dto

import "project/models"

type CreateTripRequest struct {
	Title          string `json:"title" form:"title"`
	CountryId      int    `json:"country_id" form:"country_id"`
	Accomodation   string `json:"accomodation" form:"accomodation"`
	Transportation string `json:"transportation" form:"transportation"`
	Eat            string `json:"eat" form:"eat"`
	Day            int    `json:"day" form:"day"`
	Night          int    `json:"night" form:"night"`
	DateTrip       string `json:"datetrip" form:"datetrip" vaildate:"required"`
	Price          int    `json:"price" form:"price"`
	Quota          int    `json:"quota" form:"quota"`
	Description    string `json:"description" form:"description"`
	Image          string `json:"image" form:"image"`
}

type UpdateTripRequest struct {
	Title          string `json:"title" form:"title"`
	CountryId      int    `json:"country_id" form:"country_id"`
	Accomodation   string `json:"accomodation" form:"accomodation"`
	Transportation string `json:"transportation" form:"transportation"`
	Eat            string `json:"eat" form:"eat"`
	Day            int    `json:"day" form:"day"`
	Night          int    `json:"night" form:"night"`
	DateTrip       string `json:"datetrip" form:"datetrip"`
	Price          int    `json:"price" form:"price"`
	Quota          int    `json:"quota" form:"quota"`
	Description    string `json:"description" form:"description"`
	Image          string `json:"image" form:"image"`
}

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
