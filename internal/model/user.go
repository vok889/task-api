package model

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"size:255;unique;not null"`
	Password  string `json:"password" gorm:"size:255;not null"`
	Email     string `json:"email" gorm:"size:100"`
	FirstName string `json:"first_name" gorm:"size:100"`
	LastName  string `json:"last_name" gorm:"size:100"`
}

type RequestRegister struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (u User) Exists() bool {
	return u.ID != 0 && u.Username != ""
}
