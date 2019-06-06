package model

type Item struct {
	Id          int    `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	IsValuable  bool   `gorm:"column:is_valuable" json:"-"`
	Amount      int    `gorm:"column:amount" json:"-"`
	IsAvailable bool   `gorm:"column:is_available" json:"-"`
}
