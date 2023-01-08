package models

import "time"

type Trip struct {
	Id             int             `json:"id" gorm:"primary_key:auto_increment"`
	Title          string          `json:"title" form:"title" gorm:"type: varchar(255)"`
	CountryId      int             `json:"-"`
	Country        CountryResponse `json:"country"`
	Accomodation   string          `json:"accomodation" form:"accmodation" gorm:"type: varchar(255)"`
	Transportation string          `json:"transportation" form:"transportation" gorm:"type: varchar(255)"`
	Eat            string          `json:"eat" form:"eat" gorm:"type: varchar(255)"`
	Day            int             `json:"day" form:"day" gorm:"type: int"`
	Night          int             `json:"night" form:"night" gorm:"type: int"`
	DateTrip       time.Time       `json:"datetrip"`
	Price          int             `json:"price" form:"price" gorm:"type: int"`
	Quota          int             `json:"quota" form:"quota" gorm:"type: int"`
	Description    string          `json:"description" form:"description" gorm:"type: varchar(255)"`
	Image          string          `json:"image" form:"image" gorm:"type: varchar(255)"`
}

// relation database (to transaction)
type TripResponse struct {
	Id             int             `json:"id"`
	Title          string          `json:"title"`
	CountryId      int             `json:"-"`
	Country        CountryResponse `json:"country"`
	Accomodation   string          `json:"accomodation"`
	Transportation string          `json:"transportation"`
	Eat            string          `json:"eat"`
	Day            int             `json:"day"`
	Night          int             `json:"night"`
	DateTrip       time.Time       `json:"datetrip"`
	Price          int             `json:"price"`
	Quota          int             `json:"quota"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
}

func (TripResponse) TableName() string {
	return "trips"
}
