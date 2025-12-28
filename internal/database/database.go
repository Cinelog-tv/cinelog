package database

import "gorm.io/gorm"

type Database interface {
	Connect() error
	Database() *gorm.DB
}
