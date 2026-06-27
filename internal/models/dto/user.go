package dto

type UserDTO struct {
	ID       uint    `json:"id"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Role     RoleDTO `json:"role"`
}
