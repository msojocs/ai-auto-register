package model

import "gorm.io/gorm"

type Proxy struct {
	gorm.Model
	Host     string `gorm:"not null" json:"host"`
	Port     string `gorm:"not null" json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Protocol string `gorm:"default:'http'" json:"protocol"`
	Status   string `gorm:"default:'active'" json:"status"`
}
