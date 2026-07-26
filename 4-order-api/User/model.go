package User

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone int    `json:"phone" gorm:"uniqueIndex"`
	Name  string `json:"name"`
}
