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

var myAwsInstance = awshelper.AwsInstance{}
var myDB = database.DataBase{}

type WallpaperController struct {
}

func (wc *WallpaperController) SaveWallpapers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	myDB.InitDataBase()
	resp := response.SuccessResponse{
		Status:  "Success",
		Message: "Uploaded Successfully",
	}
	category := r.FormValue("category")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error() + ", No Files found to upload"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	//handeling the file----
	file, handler, err := r.FormFile("image")
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error() + ", Error in forming file OS "
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	defer file.Close()

	//rename the file
	timestamp := time.Now().Unix()
	extension := filepath.Ext(handler.Filename)
	newFileName := fmt.Sprintf("file_%d%s", timestamp, extension)
	handler.Filename = newFileName

	//Put the file to AWS S3 instance----
	status, err := myAwsInstance.PutImageObjectToS3(file, handler, "wallpapers")

	if !status {
		resp.Status = "Failed"
		resp.Message = err.Error() + ", Unable to Save the file in Cloud"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	} else {
		//after saving into aws s3 same the sanme file name to DB too
		var wallpaper models.Wallpaper
		wallpaper.Category = category
		wallpaper.Filename = newFileName
		wallpaper.Size = handler.Size
		wallpaper.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		id, _ := utils.GenerateRandomString()
		wallpaper.WallpaperID = id
		_, err := myDB.InsertWallpaperIntoDB(wallpaper, category)
		if err != nil {
			resp.Status = "Failed"
			resp.Message = err.Error() + ", Error in Storing Image into DB"
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		response := response.SuccessResponse{
			Status:  "Success",
			Message: "Image saved successfully",
			Data:    wallpaper,
		}
		_ = json.NewEncoder(w).Encode(response)
	}
}

func (wc *WallpaperController) GetAllWallpapersByCategory(w http.ResponseWriter, r *http.Request) {
	myDB.InitDataBase()
	resp := response.SuccessResponse{
		Status:  "Success",
		Message: "Uploaded Successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	category := r.FormValue("category")
	result, err := myDB.GetWallpaperByCategory(category)
	if err != nil {
		resp.Message = err.Error()
		resp.Status = "Failed"
		resp.Data = err
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if result == nil {
		resp.Message = "No Data found for this category!"
		resp.Status = "Failed"
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = result
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)

}
