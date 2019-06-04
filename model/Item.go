package model

type Item struct {
	Id         int    `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	IsValuable bool   `gorm:"column:is_valuable"`
	Amount     int    `gorm:"column:amount"`
}

func (i Item) GetId() int {
	return i.Id
}

func (i Item) GetName() string {
	return i.Name
}
