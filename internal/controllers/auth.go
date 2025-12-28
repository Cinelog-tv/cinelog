package controllers

import (
	"fmt"
	"time"

	"github.com/Cinelog-tv/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db        *gorm.DB
	jwtSecret string
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdatePreferencesRequest struct {
	UILanguage       *string `json:"ui_language"`
	TMDBLanguage     *string `json:"tmdb_language"`
	Theme            *string `json:"theme"`
	DefaultView      *string `json:"default_view"`
	ShowAdultContent *bool   `json:"show_adult_content"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RoleResponse string

const (
	UserRoleAdmin RoleResponse = "admin"
	UserRoleUser  RoleResponse = "user"
)

type UserPreferencesResponse struct {
	UILanguage       string `json:"ui_language" gorm:"default:'en'"`
	TMDBLanguage     string `json:"tmdb_language" gorm:"default:'en-US'"`
	Theme            string `json:"theme" gorm:"default:'dark'"`
	DefaultView      string `json:"default_view" gorm:"default:'grid'"`
	ShowAdultContent bool   `json:"show_adult_content" gorm:"default:false"`
}

type UserResponse struct {
	ID          uint                    `json:"id"`
	Email       string                  `json:"email"`
	Username    string                  `json:"username"`
	FirstName   string                  `json:"first_name"`
	LastName    string                  `json:"last_name"`
	Role        RoleResponse            `json:"role"`
	Preferences UserPreferencesResponse `json:"preferences"`
	CreatedAt   time.Time               `json:"created_at"`
}

type JWTClaims struct {
	UserID uint            `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func RegisterAuthController(api *echo.Group, db *gorm.DB, jwtSecret string) *AuthController {
	a := &AuthController{
		db:        db,
		jwtSecret: jwtSecret,
	}

	// Public routes
	api.POST("/auth/register", a.registerController)
	api.POST("/auth/login", a.loginController)

	// Protected routes (require authentication)
	api.GET("/auth/me", a.getCurrentUserController, a.authMiddleware)
	api.POST("/auth/logout", a.logoutController)
	api.PUT("/auth/preferences", a.updatePreferencesController, a.authMiddleware)
	api.PUT("/auth/profile", a.updateProfileController, a.authMiddleware)

	// Admin only routes
	api.POST("/admin/users", a.createUserController, a.authMiddleware, a.adminMiddleware)
	api.GET("/admin/users", a.listUsersController, a.authMiddleware, a.adminMiddleware)
	api.GET("/admin/users/:id", a.getUserController, a.authMiddleware, a.adminMiddleware)
	api.DELETE("/admin/users/:id", a.deleteUserController, a.authMiddleware, a.adminMiddleware)
	api.PUT("/admin/users/:id/role", a.updateUserRoleController, a.authMiddleware, a.adminMiddleware)

	return a
}

func (a *AuthController) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(401, map[string]string{"error": "Missing authorization token"})
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(401, map[string]string{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.JSON(401, map[string]string{"error": "Invalid token claims"})
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		return next(c)
	}
}

func (a *AuthController) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("user_role").(models.UserRole)
		if !ok || role != models.UserRoleAdmin {
			return c.JSON(403, map[string]string{"error": "Forbidden - Admin access required"})
		}
		return next(c)
	}
}

// ShowAccount godoc
// @Summary      Create user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/register [post]
func (a *AuthController) registerController(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	var adminCount int64
	a.db.Model(&models.User{}).Where("role = ?", models.UserRoleAdmin).Count(&adminCount)

	if adminCount > 0 {
		return c.JSON(403, map[string]string{"error": "Public registration is disabled. Please contact an administrator."})
	}

	var existingUser models.User
	if err := a.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(400, map[string]string{"error": "Email already exists"})
	}

	if err := a.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.JSON(400, map[string]string{"error": "Username already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to hash password"})
	}

	// First user becomes admin
	role := models.UserRoleAdmin

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         role,
	}

	if err := a.db.Create(&user).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create user"})
	}

	preferences := models.UserPreferences{
		UserID:           user.ID,
		UILanguage:       "en",
		TMDBLanguage:     "en-US",
		Theme:            "dark",
		DefaultView:      "grid",
		ShowAdultContent: false,
	}

	if err := a.db.Create(&preferences).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create user preferences"})
	}

	token, err := a.generateJWT(user.ID, user.Role)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	a.db.Preload("Preferences").First(&user, user.ID)

	return c.JSON(201, AuthResponse{
		Token: token,
		User:  mapUserToResponse(user),
	})
}

func (a *AuthController) loginController(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	var user models.User
	if err := a.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	token, err := a.generateJWT(user.ID, user.Role)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	a.db.Preload("Preferences").First(&user, user.ID)

	return c.JSON(200, AuthResponse{
		Token: token,
		User:  mapUserToResponse(user),
	})
}

func (a *AuthController) getCurrentUserController(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var user models.User
	if err := a.db.Preload("Preferences").First(&user, userID).Error; err != nil {
		return c.JSON(404, map[string]string{"error": "User not found"})
	}

	return c.JSON(200, mapUserToResponse(user))
}

func (a *AuthController) logoutController(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "Logged out successfully"})
}

func (a *AuthController) updatePreferencesController(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	var preferences models.UserPreferences
	result := a.db.Where("user_id = ?", userID).First(&preferences)

	if result.Error != nil {
		preferences = models.UserPreferences{
			UserID:           userID,
			UILanguage:       "en",
			TMDBLanguage:     "en-US",
			Theme:            "dark",
			DefaultView:      "grid",
			ShowAdultContent: false,
		}
		if err := a.db.Create(&preferences).Error; err != nil {
			return c.JSON(500, map[string]string{"error": "Failed to create preferences"})
		}
	}

	if req.UILanguage != nil {
		preferences.UILanguage = *req.UILanguage
	}
	if req.TMDBLanguage != nil {
		preferences.TMDBLanguage = *req.TMDBLanguage
	}
	if req.Theme != nil {
		preferences.Theme = *req.Theme
	}
	if req.DefaultView != nil {
		preferences.DefaultView = *req.DefaultView
	}
	if req.ShowAdultContent != nil {
		preferences.ShowAdultContent = *req.ShowAdultContent
	}

	if err := a.db.Save(&preferences).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to update preferences"})
	}

	return c.JSON(200, preferences)
}

func (a *AuthController) updateProfileController(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	var user models.User
	if err := a.db.First(&user, userID).Error; err != nil {
		return c.JSON(404, map[string]string{"error": "User not found"})
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Username != "" {
		var existingUser models.User
		if err := a.db.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
			return c.JSON(400, map[string]string{"error": "Username already taken"})
		}
		user.Username = req.Username
	}

	if err := a.db.Save(&user).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to update profile"})
	}

	a.db.Preload("Preferences").First(&user, userID)
	return c.JSON(200, mapUserToResponse(user))
}

func (a *AuthController) createUserController(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	var existingUser models.User
	if err := a.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(400, map[string]string{"error": "Email already exists"})
	}

	if err := a.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.JSON(400, map[string]string{"error": "Username already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to hash password"})
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         models.UserRoleUser,
	}

	if err := a.db.Create(&user).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create user"})
	}

	preferences := models.UserPreferences{
		UserID:           user.ID,
		UILanguage:       "en",
		TMDBLanguage:     "en-US",
		Theme:            "dark",
		DefaultView:      "grid",
		ShowAdultContent: false,
	}

	if err := a.db.Create(&preferences).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create user preferences"})
	}

	a.db.Preload("Preferences").First(&user, user.ID)

	return c.JSON(201, mapUserToResponse(user))
}

func (a *AuthController) listUsersController(c echo.Context) error {
	var users []models.User
	if err := a.db.Preload("Preferences").Find(&users).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch users"})
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = mapUserToResponse(user)
	}

	return c.JSON(200, userResponses)
}

func (a *AuthController) getUserController(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := a.db.Preload("Preferences").First(&user, id).Error; err != nil {
		return c.JSON(404, map[string]string{"error": "User not found"})
	}

	return c.JSON(200, mapUserToResponse(user))
}

func (a *AuthController) deleteUserController(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(uint)

	var idUint uint
	fmt.Sscanf(id, "%d", &idUint)
	if idUint == userID {
		return c.JSON(400, map[string]string{"error": "Cannot delete your own account"})
	}

	if err := a.db.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to delete user"})
	}

	return c.JSON(200, map[string]string{"message": "User deleted successfully"})
}

func (a *AuthController) updateUserRoleController(c echo.Context) error {
	id := c.Param("id")

	var req struct {
		Role models.UserRole `json:"role" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if req.Role != models.UserRoleAdmin && req.Role != models.UserRoleUser {
		return c.JSON(400, map[string]string{"error": "Invalid role"})
	}

	var user models.User
	if err := a.db.First(&user, id).Error; err != nil {
		return c.JSON(404, map[string]string{"error": "User not found"})
	}

	user.Role = req.Role
	if err := a.db.Save(&user).Error; err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to update role"})
	}

	a.db.Preload("Preferences").First(&user, id)
	return c.JSON(200, mapUserToResponse(user))
}

func (a *AuthController) generateJWT(userID uint, role models.UserRole) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func mapUserToResponse(user models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      RoleResponse(user.Role),
		Preferences: UserPreferencesResponse{
			TMDBLanguage:     user.Preferences.TMDBLanguage,
			UILanguage:       user.Preferences.UILanguage,
			Theme:            user.Preferences.Theme,
			DefaultView:      user.Preferences.DefaultView,
			ShowAdultContent: user.Preferences.ShowAdultContent,
		},
		CreatedAt: user.CreatedAt,
	}
}
