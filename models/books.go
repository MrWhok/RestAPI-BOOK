package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `json:"id" gorm:"primarykey;autoIncrement"`
	Author    *string `json:"author"` //why using *string? because it can be null (optional data)
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{}) //what is automigrate? it will create the table if it doesnt exist and if it exist it will update the table. If it doesnt use automigrate, we must do it manually like "ALTER TABLE"
	return err
}
