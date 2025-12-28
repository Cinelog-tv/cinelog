package dto

type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID          uint                     `json:"id"`
	Email       string                   `json:"email"`
	Username    string                   `json:"username"`
	FirstName   string                   `json:"first_name"`
	LastName    string                   `json:"last_name"`
	Role        string                   `json:"role"`
	Preferences *UserPreferencesResponse `json:"preferences,omitempty"`
	CreatedAt   string                   `json:"created_at"`
}

type UserPreferencesResponse struct {
	ID               uint   `json:"id"`
	UILanguage       string `json:"ui_language"`
	TMDBLanguage     string `json:"tmdb_language"`
	Theme            string `json:"theme"`
	DefaultView      string `json:"default_view"`
	ShowAdultContent bool   `json:"show_adult_content"`
}

type UpdatePreferencesRequest struct {
	UILanguage       string `json:"ui_language"`
	TMDBLanguage     string `json:"tmdb_language"`
	Theme            string `json:"theme"`
	DefaultView      string `json:"default_view"`
	ShowAdultContent *bool  `json:"show_adult_content"`
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}
