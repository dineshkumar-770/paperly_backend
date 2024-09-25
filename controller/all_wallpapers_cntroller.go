package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"mongo_api/helpers"
	"mongo_api/models"
	"mongo_api/response"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

// RetrieveS3FileInstance
 
func RetrieveAllImageFromBucket(w http.ResponseWriter, r *http.Request){
 
	w.Header().Set("Content-Type", "application/json")
	resp := response.SuccessResponse{
		Status: "Failed",
	}
	var allWallpapers []models.Wallpaper

	_ = godotenv.Load(".env") 
	awsBucket := os.Getenv("BUCKETNAME")
	svc := helpers.GetAllFilesFromBucket()
	input := &s3.ListObjectsV2Input{
        Bucket: aws.String(awsBucket),
    }
	result, err := svc.ListObjectsV2(input)
	if err !=nil{
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

	for _,item := range result.Contents{
		var wallpaper models.Wallpaper
		if (strings.Contains(*item.Key,"category")){
			fmt.Println("Contains Category Images : ", strings.Contains(*item.Key,"category"))
		}else{
			if isImage(*item.Key){ 
				req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(awsBucket),
					Key:    aws.String(*item.Key),
				})
	
				urlStr,err :=req.Presign(1 * time.Hour)
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

	if allWallpapers == nil{
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