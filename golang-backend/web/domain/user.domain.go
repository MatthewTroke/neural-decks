package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	UpsertUserByID(user *User) (*User, error)
}

type User struct {
	ID            string `gorm:"primaryKey;type:varchar(255)"`
	Name          string `gorm:"column:name"`
	Email         string `gorm:"unique;column:email"`
	EmailVerified bool   `gorm:"column:email_verified"`
	Image         string `gorm:"column:image"`
	Provider      string `gorm:"column:provider"`
}
