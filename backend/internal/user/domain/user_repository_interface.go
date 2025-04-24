package domain

type IUserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}
