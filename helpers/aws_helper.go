package helpers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mongo_api/utils"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func AddImageToS3(file multipart.File, fileHeader *multipart.FileHeader, filePathS3 string) (status bool, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading Environments", err)
		return false, err
	}
	envVars, errEnv := utils.GetEnvVariables()
	if envVars.BucketName == "" {
		return false, errEnv
	}
	awsRegion := envVars.AWSRegion
	awsBucket := envVars.BucketName
	s3FolderPath := filePathS3 + "/"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Fatal("Error creating session: ", err)
		return false, err
	}

	_, err = io.ReadAll(file)
	if err != nil {
		log.Fatal("Error in reading file: ", err)
		return false, err
	}
	fileSize, err := file.Seek(0, io.SeekEnd) // Move to end to get the size
	if err != nil {
		log.Fatal("Error getting file size: ", err)
		return false, err
	}
	fmt.Printf("File size: %d bytes\n", fileSize)

	// Reset the file pointer to the start
	_, err = file.Seek(0, io.SeekStart) // Reset pointer to start for upload
	if err != nil {
		log.Fatal("Error resetting file pointer: ", err)
		return false, err
	}

	defer file.Close()

	svc := s3.New(sess)

	input := s3.PutObjectInput{
		Bucket:      aws.String(awsBucket),
		Key:         aws.String(path.Join(s3FolderPath, fileHeader.Filename)),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	}

	_, err = svc.PutObject(&input)
	if err != nil {
		log.Fatal("Error in Uploading file: ", err)
		return false, err
	}

	fmt.Printf("File Uploaded Successfully!")
	return true, err
}

func GetAllFilesFromBucket() *s3.S3 {
	_ = godotenv.Load(".env")
	envVars, _ := utils.GetEnvVariables()
	if envVars.BucketName == "" {
		return nil
	}
	awsRegion := envVars.AWSRegion
	// awsBucket := os.Getenv("BUCKETNAME")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Fatal("Error creating session: ", err)
	}

	if err != nil {
		log.Fatal("Error in reading file: ", err)
	}

	svc := s3.New(sess)
	return svc
}

// func RetrieveAllImageFromBucket(){
// 	err = godotenv.Load(".env")
// 	if err !=nil{
// 		log.Fatal("Error in loading Environments",err)
// 		// return false,err
// 	}
// 	awsRegion := os.Getenv("AWSREGION")
// 	awsBucket := os.Getenv("BUCKETNAME")

// 	sess, err :=  session.NewSession(&aws.Config{
// 		Region: aws.String(awsRegion),
// 	})
// 	if err != nil{
// 		log.Fatal("Error creating session: ",err)
// 		// return false ,err
// 	}
// }
