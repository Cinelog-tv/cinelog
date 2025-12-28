package models

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypeMovie   MediaType = "movie"
	MediaTypeTV      MediaType = "tv"
	MediaTypeEpisode MediaType = "episode"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type WatchedMedia struct {
	gorm.Model
	UserID     uint      `json:"user_id" gorm:"uniqueIndex:idx_watched_media;not null"`
	TmdbID     int       `json:"tmdb_id" gorm:"uniqueIndex:idx_watched_media"`
	MediaType  MediaType `json:"media_type" gorm:"uniqueIndex:idx_watched_media"`
	SeasonNum  int       `json:"season_number"`  // For episodes only
	EpisodeNum int       `json:"episode_number"` // For episodes only
	Watched    bool      `json:"watched" gorm:"default:false"`
	WatchCount int       `json:"watch_count" gorm:"default:0"`
	User       User      `json:"-" gorm:"foreignKey:UserID"`
}

func (WatchedMedia) TableName() string {
	return "watched_media"
}

type PhysicalMedia struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	TmdbID     int            `json:"tmdb_id"`
	MediaType  MediaType      `json:"media_type"`
	SeasonNum  int            `json:"season_number"` // For TV seasons only
	Format     string         `json:"format" gorm:"not null"`
	Edition    string         `json:"edition"`
	Price      float64        `json:"price"`
	Store      string         `json:"store"`
	PurchaseAt *gorm.DeletedAt `json:"purchase_at" gorm:"index"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
}

func (PhysicalMedia) TableName() string {
	return "physical_media"
}

type User struct {
	gorm.Model
	Email        string       `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string       `json:"-" gorm:"not null"`
	Username     string       `json:"username" gorm:"uniqueIndex;not null"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Role         UserRole     `json:"role" gorm:"default:'user';not null"`
	Preferences  *UserPreferences `json:"preferences" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

type UserPreferences struct {
	gorm.Model
	UserID           uint   `json:"user_id" gorm:"uniqueIndex;not null"`
	UILanguage       string `json:"ui_language" gorm:"default:'en'"`
	TMDBLanguage     string `json:"tmdb_language" gorm:"default:'en-US'"`
	Theme            string `json:"theme" gorm:"default:'dark'"`
	DefaultView      string `json:"default_view" gorm:"default:'grid'"`
	ShowAdultContent bool   `json:"show_adult_content" gorm:"default:false"`
	User             User   `json:"-" gorm:"foreignKey:UserID"`
}

func (UserPreferences) TableName() string {
	return "user_preferences"
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &UserPreferences{}, &WatchedMedia{}, &PhysicalMedia{})
}
