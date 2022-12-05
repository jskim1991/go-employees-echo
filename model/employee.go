package model

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name   string `gorm:"type:varchar(255)"`
	Salary string `gorm:"type:varchar(255)"`
	Age    int    `gorm:"type:uint"`
}
