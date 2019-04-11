package controller

import (
	"agit.com/smartdashboard-backend/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB = config.LoadDB()
