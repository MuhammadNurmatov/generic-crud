package example

type UserExample struct {
	ID    int64  `gorm:"primary_key"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}
