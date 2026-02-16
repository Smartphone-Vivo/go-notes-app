package service

type Note struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Text
}
