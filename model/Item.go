package model

type Item struct {
	Id          int    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	IsAvailable bool   `gorm:"column:is_available" json:"-"`
}
