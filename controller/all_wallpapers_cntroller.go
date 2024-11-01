package controller

import (
	"encoding/json"
	"fmt"
	"log"
	awshelper "mongo_api/aws_helper"
	"mongo_api/database"
	"mongo_api/helpers"
	"mongo_api/models"
	"mongo_api/response"
	"mongo_api/utils"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var dbInstnce = database.DataBase{}
	var awsInstnce = awshelper.AwsInstance{}

func DeleteCategoryImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbInstnce.InitDataBase()
	awsInstnce.AwsInit()
	resp := response.SuccessResponse{
		Status: "Failed",
	}

	var wallpaperObject models.Wallpaper
	err := json.NewDecoder(r.Body).Decode(&wallpaperObject)
	if err != nil {
		w.WriteHeader(401)
		resp.Message = "Unable to parse the wallpaper object. kindly provider valid wallpaper to delete!"
		resp.Data = nil
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Println("decoded wall object: ", wallpaperObject)

	_, err1 := awsInstnce.DeleteFileFromS3(wallpaperObject.Filename)
	if err1 != nil {
		w.WriteHeader(500)
		resp.Message = "Internal Server Error from Cloud Storage. Please try again later!"
		resp.Data = nil
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		status, _ := dbInstnce.DeleteOneImage(wallpaperObject)
		if !status {
			w.WriteHeader(500)
			resp.Message = "Internal Server Error from Data storage. Please try again later!"
			resp.Data = nil
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(200)
		resp.Status = "Success"
		resp.Message = "Image deleted successfully!"
		json.NewEncoder(w).Encode(resp)
		return
	}

	

}

// RetrieveS3FileInstance

func RetrieveAllImageFromBucket(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	resp := response.SuccessResponse{
		Status: "Failed",
	}
	var allWallpapers []models.Wallpaper
	envVars, errEnv := utils.GetEnvVariables()
	if envVars.BucketName == "" {
		resp.Status = "Failed"
		resp.Message = errEnv.Error() + ", No environment found"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	awsBucket := envVars.BucketName
	folderS3 := envVars.BucketFolderName
	svc := helpers.GetAllFilesFromBucket()
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(awsBucket),
		Prefix: aws.String(folderS3),
	}
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		log.Println("Error in getting images from S3:-- ", err)
		resp.Message = "Error in fetching Images from cloud"
		resp.Data = err
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(resp)
		return
	}
	sort.Slice(result.Contents, func(i, j int) bool {
		return result.Contents[i].LastModified.After(*result.Contents[j].LastModified)
	})

	for _, item := range result.Contents {
		var wallpaper models.Wallpaper
		if strings.Contains(*item.Key, "category") {
			fmt.Println("Contains Category Images : ", strings.Contains(*item.Key, "category"))
		} else {
			if isImage(*item.Key) {
				req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(awsBucket),
					Key:    aws.String(*item.Key),
				})

				urlStr, err := req.Presign(24 * time.Hour)
				if err != nil {
					log.Println("error generating presigned URL for: --- ", err)
					resp.Message = "Error in fetching Images from cloud"
					resp.Data = err
					w.WriteHeader(500)
					json.NewEncoder(w).Encode(resp)
					break
				}

				wallpaper.Size = *item.Size
				wallpaper.Filename = urlStr
				wallpaper.Category = ""

				allWallpapers = append(allWallpapers, wallpaper)
			}
		}

	}

	if allWallpapers == nil {
		w.WriteHeader(403)
		resp.Message = "No images are available to display at the moment."
		resp.Status = "Failed"
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp.Status = "Success"
	resp.Message = "Successfully fetched all Images"
	resp.Data = allWallpapers
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func isImage(key string) bool {
	lowerKey := strings.ToLower(key)
	return strings.HasSuffix(lowerKey, ".jpg") || strings.HasSuffix(lowerKey, ".jpeg") ||
		strings.HasSuffix(lowerKey, ".png") || strings.HasSuffix(lowerKey, ".gif")
}
