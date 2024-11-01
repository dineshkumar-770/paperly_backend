package main

import (

	// "io"
	"log"
	awshelper "mongo_api/aws_helper"
	"mongo_api/controller"
	"mongo_api/database"
	"net/http"

	"github.com/gorilla/mux"
)

var myDB = database.DataBase{}
var awsInstance = awshelper.AwsInstance{}
var wallpaperController = controller.WallpaperController{}
var wallpaperCategoriesCont = controller.WallCategories{}
var deviceInfoAndFCMController = controller.DeviceInfoAndFCM{}
var notiController = controller.NotificationController{}

func main() {
	awsInstance.AwsInit()
	myDB.InitDataBase()
	r := mux.NewRouter()
	r.HandleFunc("/add_wallpaper", wallpaperController.SaveWallpapers).Methods("POST")
	r.HandleFunc("/get_all_wall_by_category", wallpaperController.GetAllWallpapersByCategory).Methods("POST")
	r.HandleFunc("/add_category", wallpaperCategoriesCont.AddWallpaperCategories).Methods("POST")
	r.HandleFunc("/get_all_images", controller.RetrieveAllImageFromBucket).Methods("GET")
	r.HandleFunc("/get_all_categories", wallpaperCategoriesCont.GetAllCategories).Methods("GET")
	r.HandleFunc("/save_fcm_token", deviceInfoAndFCMController.SaveUSerDeviceInfoWithFCM).Methods("POST")
	r.HandleFunc("/send_simple_notification",notiController.SendSimpleNotification).Methods("POST")
	r.HandleFunc("/delete_category_image",controller.DeleteCategoryImage).Methods("POST")
	log.Fatal(http.ListenAndServe(":4400", r))
}
