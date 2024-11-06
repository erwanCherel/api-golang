package models

import (
	"time"
)

// User model representing the user table in the database
type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Pseudo    string    `json:"pseudo,omitempty"`
	Email     string    `gorm:"unique;not null" json:"email,omitempty"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Videos    []Video   `gorm:"foreignKey:UserID" json:"-"`
	Comments  []Comment `gorm:"foreignKey:UserID" json:"-"`
	Tokens    []Token   `gorm:"foreignKey:UserID" json:"-"`
}

// Video model representing the video table in the database
type Video struct {
	ID        int           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string        `gorm:"not null" json:"-"`
	Duration  int           `json:"-"`
	UserID    int           `gorm:"not null" json:"-"`
	User      User          `gorm:"foreignKey:UserID" json:"user"`
	Source    string        `gorm:"not null" json:"source"`
	CreatedAt time.Time     `json:"created_at"`
	Views     int           `json:"views"`
	Enabled   bool          `json:"enabled"`
	Formats   []VideoFormat `gorm:"foreignKey:VideoID" json:"format"`
	Comments  []Comment     `gorm:"foreignKey:VideoID" json:"-"`
}

// VideoFormat model representing the video_format table in the database
type VideoFormat struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"-"`
	Code    string `gorm:"not null" json:"-"`
	URI     string `gorm:"not null" json:"-"`
	VideoID int    `gorm:"not null" json:"-"`
}

// Token model representing the token table in the database
type Token struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"-"`
	Code      string    `gorm:"unique;not null" json:"token"`
	ExpiredAt time.Time `json:"-"`
	UserID    int       `gorm:"not null" json:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// Comment model representing the comment table in the database
type Comment struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Body    string `gorm:"type:longtext" json:"body"`
	UserID  int    `gorm:"not null" json:"-"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	VideoID int    `gorm:"not null" json:"-"`
}

// Custom format mapping for the video format resolution URIs
type VideoFormatMap struct {
	Resolution1080 string `json:"1080,omitempty"`
	Resolution720  string `json:"720,omitempty"`
	Resolution480  string `json:"480,omitempty"`
	Resolution360  string `json:"360,omitempty"`
	Resolution240  string `json:"240,omitempty"`
	Resolution144  string `json:"144,omitempty"`
}

// To map video formats by resolution, create a helper function
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
