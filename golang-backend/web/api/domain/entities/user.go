package entities

type User struct {
	ID            string `gorm:"primaryKey;type:varchar(255)"`
	Name          string `gorm:"column:name"`
	Email         string `gorm:"unique;column:email"`
	EmailVerified bool   `gorm:"column:email_verified"`
	Image         string `gorm:"column:image"`
	Provider      string `gorm:"column:provider"`
}
