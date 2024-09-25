package database

import (
	"context"
	"mongo_api/helpers"
	"mongo_api/models"
	"time"

	// "errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
	Client *mongo.Client
}

var dbName string = "wallpapers"

func (d *DataBase) InitDataBase() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		return nil
	}
	dburl := os.Getenv("DATABASEURL")
	fmt.Println(dburl)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dburl).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	d.Client = client
	fmt.Println("Database connected successfully")
	return client
}

func (d *DataBase) InsertWallpaperIntoDB(wallpaper models.Wallpaper, category string) (*mongo.InsertOneResult, error) {
	if d.Client == nil {
		return nil, fmt.Errorf("database client is not initialized")
	}
	collection := d.Client.Database("wallpapers").Collection(category)
	result, err := collection.InsertOne(context.TODO(), wallpaper)
	return result, err
}

func (d *DataBase) GetWallpaperByCategory(category string) ([]models.Wallpaper, error) {
	_ = godotenv.Load(".env")
	// awsRegion := os.Getenv("AWSREGION")
	awsBucket := os.Getenv("BUCKETNAME")
	if d.Client == nil {
		return nil, fmt.Errorf("database client is not initialized")
	}

	collection := d.Client.Database(dbName).Collection(category)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("unable to found wallpapers")
		return nil, err
	}

	defer cursor.Close(context.Background())
	var wallpapers []models.Wallpaper

	for cursor.Next(context.TODO()) {
		var wallpaper models.Wallpaper
		cursor.Decode(&wallpaper)
		key := fmt.Sprintf("wallpapers/%s", wallpaper.Filename)
		svc := helpers.GetAllFilesFromBucket()
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(awsBucket),
			Key:    aws.String(key),
		})

		urlStr, err := req.Presign(1 * time.Hour)
		if err != nil {
			log.Printf("Error generating presigned URL for %s: %v", wallpaper.Filename, err)
			continue
		}

		wallpaper.Filename = urlStr
		wallpapers = append(wallpapers, wallpaper)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return wallpapers, nil
}

func (d *DataBase) GetAllCategoriesList() ([]models.WallPaperCategories, error) {
	_ = godotenv.Load(".env")
	awsBucket := os.Getenv("BUCKETNAME")
	if d.Client == nil {
		return nil, fmt.Errorf("database client is not initialized")
	}

	collection := d.Client.Database(dbName).Collection("wallpapers_categories")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("unable to found wallpapers")
		return nil, err
	}

	defer cursor.Close(context.Background())

	var allCategories []models.WallPaperCategories

	for cursor.Next(context.TODO()) {
		var category models.WallPaperCategories
		cursor.Decode(&category)
		key := fmt.Sprintf("wallpapers/%s", category.CategoryImage)
		svc := helpers.GetAllFilesFromBucket()
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(awsBucket),
			Key:    aws.String(key),
		})

		urlStr, err := req.Presign(1 * time.Hour)
		if err != nil {
			log.Printf("Error generating presigned URL for %s: %v", category.CategoryImage, err)
			continue
		}

		category.CategoryImage = urlStr
		allCategories = append(allCategories, category)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return allCategories, nil
}

func (d *DataBase) AddCategories(category models.WallPaperCategories) (bool, error) {
	if d.Client == nil {
		return false, fmt.Errorf("database client is not initialized")
	}
	collection := d.Client.Database("wallpapers").Collection("wallpapers_categories")
	_, err := collection.InsertOne(context.TODO(), category)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (d *DataBase) SaveDeviceInfo(deviceInfo models.DeviceInfo) (bool, error) {
	if d.Client == nil {
		return false, fmt.Errorf("database client is not initialized")
	}

	collection := d.Client.Database("wallpapers").Collection("device_information")
	_, err := collection.InsertOne(context.TODO(), deviceInfo)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *DataBase) FindOneCategory(categoryName string) (bool, error) {
	if d.Client == nil {
		return true, fmt.Errorf("database client is not initialized")
	}
	var categoryData models.WallPaperCategories
	collection := d.Client.Database("wallpapers").Collection("wallpapers_categories")
	filter := bson.M{"category_name": categoryName}
	err := collection.FindOne(context.TODO(), filter).Decode(&categoryData)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	log.Println(err)

	return true, err

}
