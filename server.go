package main

import "gorm.io/gorm"

type Server struct {
	Db *gorm.DB
}
