package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	BucketName       string
	BucketFolderName string
	DatabaseURL      string
	AWSRegion        string
	AWSAccessKey     string
	AWSSecretKey     string
	DatabaseName     string
}

func GetEnvVariables() (EnvVariables, error) {
	var e EnvVariables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return e, err
	}

	databaseURL := os.Getenv("PRODUCTIONDB")
	awsRegion := os.Getenv("AWSREGION")
	bucketName := os.Getenv("BUCKETNAME")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	databaseName := os.Getenv("PRODDBNAME")
	bucketFolderName := os.Getenv("PRODUCTIONBUCKETFOLDER")

	e.AWSAccessKey = awsAccessKey
	e.AWSRegion = awsRegion
	e.AWSSecretKey = awsSecretKey
	e.DatabaseName = databaseName
	e.BucketName = bucketName
	e.DatabaseURL = databaseURL
	e.BucketFolderName = bucketFolderName

	return e, nil
}
