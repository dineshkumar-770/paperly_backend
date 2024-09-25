package controller

import (
	"encoding/json"
	"fmt"
	awshelper "mongo_api/aws_helper"
	"mongo_api/database"
	"mongo_api/models"
	"mongo_api/response"
	"mongo_api/utils"
	"net/http"
	"path/filepath"
	"time"
)

type WallCategories struct {
	Wallpapercategories models.WallPaperCategories `json:"wallpaper_categories" bson:"wallpaper_categories"`
}

var myDatabase = database.DataBase{}
var categoryAwsInstance = awshelper.AwsInstance{}

func (wal *WallCategories) AddWallpaperCategories(w http.ResponseWriter, r *http.Request) {
	resp := response.SuccessResponse{
		Status:  "Success",
		Message: "category resigtered successfully",
		Data:    nil,
	}
	var wallCategory models.WallPaperCategories
	w.Header().Set("Content-Type", "multipart/form-data")
	wallCategoryName := r.FormValue("category")
	myDatabase.InitDataBase()

	if wallCategoryName == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp.Status = "failed"
		resp.Message = "category name is required to register it"
		json.NewEncoder(w).Encode(resp)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error() + ", No Files found to upload"
		w.WriteHeader(403)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error() + ", Error in forming file OS "
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	defer file.Close()

	timestamp := time.Now().Unix()
	extension := filepath.Ext(handler.Filename)
	newFileName := fmt.Sprintf("category_%d%s", timestamp, extension)
	handler.Filename = newFileName

	wallID, _ := utils.GenerateRandomString()
	wallCategory.CategoryId = wallID
	wallCategory.TotalImages = 0
	wallCategory.CategoryName = wallCategoryName
	wallCategory.CategoryImage = newFileName

	getPreexistCategory, err := myDatabase.FindOneCategory(wallCategory.CategoryName)

	if !getPreexistCategory {
		status, errr := categoryAwsInstance.PutImageObjectToS3(file, handler, "wallpapers")
		if !status {
			resp.Status = "Failed"
			resp.Message = errr.Error() + ", Unable to Save the file in Cloud"
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		_, err := myDatabase.AddCategories(wallCategory)
		if err != nil {
			resp.Status = "failed"
			resp.Message = "Unable to register category at this time, please try again later!"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		resp.Status = "failed"
		resp.Message = "cannot create category with the name already exists!"
		resp.Data = err
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func (wal *WallCategories) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	myDatabase.InitDataBase()
	resp := response.SuccessResponse{
		Status:  "Failed",
		Message: "",
	}

	result, err := myDatabase.GetAllCategoriesList()
	if err != nil {
		resp.Message = err.Error()
		resp.Status = "Failed"
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if result == nil {
		w.WriteHeader(403)
		resp.Status = "failed"
		resp.Message = "No categories are available at this moment!"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Status = "Success"
	resp.Message = "All Categories fetched successfully"
	resp.Data = result
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// func (wal *WallCategories) GetAllCategories(w http.ResponseWriter, r *http.Request) {

// 	category := r.FormValue("category")

// 	resp := response.SuccessResponse{
// 		Status:  "Failed",
// 		Message: "Unable to generate token",
// 	}

// 	token, privateKey, err := utils.CreateJWTToken(category)
// 	if err != nil {
// 		w.WriteHeader(403)
// 		resp.Data = token
// 		json.NewEncoder(w).Encode(resp)
// 		return
// 	}

// 	resp.Data = map[string]interface{}{
// 		"token": token,
// 		"key": privateKey.PublicKey,
// 	}

// 	resp.Message = "Token Created succesfully"
// 	resp.Status = "Success"
// 	log.Println(&privateKey.PublicKey)
// 	json.NewEncoder(w).Encode(resp)

// }
