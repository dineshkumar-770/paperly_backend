package models


type Wallpaper struct {
	Filename    string `json:"filename" bson:"filename"`
	Size        int64 `json:"size" bson:"size"`
	Timestamp   string `json:"timestamp" bson:"timestamp"`
	WallpaperID string `json:"wallpaper_id" bson:"wallpaper_id"`
	Category    string `json:"category" bson:"category"`
}