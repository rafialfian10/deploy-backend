package models

type Image struct {
	Id     int    `json:"id"`
	Name   string `json:"file_name" gorm:"type:varchar(255)"`
	TripID int
	Trip   TripResponse
}

type ImageResponse struct {
	Id     int    `json:"-"`
	Name   string `json:"file_name"`
	TripID int    `json:"-"`
}

func (ImageResponse) TableName() string {
	return "images"
}
