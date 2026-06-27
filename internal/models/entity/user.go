package entity

type UserEntity struct {
	ID       uint
	Email    string
	Username string
	Password string
	Role     *RoleEntity
}
