package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDatabase struct {
	db *gorm.DB
}

func (s *SqliteDatabase) Connect() error {
	db, err := gorm.Open(sqlite.Open("cinelog.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SqliteDatabase) Database() *gorm.DB {
	return s.db
}
