package models

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Pseudo    string    `json:"pseudo,omitempty"`
	Email     string    `gorm:"unique;not null" json:"email,omitempty"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Videos    []Video   `gorm:"foreignKey:UserID" json:"-"`
	Comments  []Comment `gorm:"foreignKey:UserID" json:"-"`
	Tokens    []Token   `gorm:"foreignKey:UserID" json:"-"`
}

type Video struct {
	ID        int           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string        `gorm:"not null" json:"name"`
	UserID    int           `gorm:"not null" json:"user_id"`
	Source    string        `gorm:"not null" json:"source"`
	CreatedAt time.Time     `json:"created_at"`
	Views     int           `json:"views"`
	Enabled   bool          `json:"enabled"`
	Formats   []VideoFormat `gorm:"foreignKey:VideoID" json:"formats"`
	Comments  []Comment     `gorm:"foreignKey:VideoID;constraint:OnDelete:CASCADE;" json:"comments"`
}

type VideoFormat struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code    string `gorm:"not null" json:"code"`
	URI     string `gorm:"not null" json:"uri"`
	VideoID int    `gorm:"not null" json:"video_id"`
}

type Token struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"unique;not null" json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	UserID    int       `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

type Comment struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Body    string `gorm:"type:longtext" json:"body"`
	UserID  int    `gorm:"not null" json:"user_id"`
	VideoID int    `gorm:"not null;constraint:OnDelete:CASCADE;" json:"video_id"`
}

type VideoFormatMap struct {
	Resolution1080 string `json:"1080,omitempty"`
	Resolution720  string `json:"720,omitempty"`
	Resolution480  string `json:"480,omitempty"`
	Resolution360  string `json:"360,omitempty"`
	Resolution240  string `json:"240,omitempty"`
	Resolution144  string `json:"144,omitempty"`
}

func (v *Video) FormatMap() VideoFormatMap {
	formatMap := VideoFormatMap{}
	for _, format := range v.Formats {
		switch format.Code {
		case "1080":
			formatMap.Resolution1080 = format.URI
		case "720":
			formatMap.Resolution720 = format.URI
		case "480":
			formatMap.Resolution480 = format.URI
		case "360":
			formatMap.Resolution360 = format.URI
		case "240":
			formatMap.Resolution240 = format.URI
		case "144":
			formatMap.Resolution144 = format.URI
		}
	}
	return formatMap
}
