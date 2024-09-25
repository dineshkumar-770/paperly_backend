package models

type WallPaperCategories struct{
	CategoryName string `json:"category_name" bson:"category_name"`
	CategoryId string `json:"category_id" bson:"category_id"`
	TotalImages int `json:"total_images" bson:"total_images"`
	CategoryImage string `json:"category_image" bson:"category_image"`
}