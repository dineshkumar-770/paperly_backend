package awshelper

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mongo_api/utils"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsInstance struct {
	BucketName string
	AwsRegion  string
	MyError    error
}

func (a *AwsInstance) AwsInit() (*session.Session, error) {

	envVars, errEnv := utils.GetEnvVariables()
	if envVars.BucketName == "" {
		return nil, errEnv
	}

	a.AwsRegion = envVars.AWSRegion
	a.BucketName = envVars.BucketName

	aws_sesson, err := session.NewSession(
		&aws.Config{
			Region: aws.String(a.BucketName),
		},
	)

	if err != nil {
		a.MyError = err
		return nil, a.MyError
	}

	return aws_sesson, a.MyError
}

func (a *AwsInstance) DeleteFileFromS3(fileKey string) (bool, error) {
	envVars, err := utils.GetEnvVariables()

	if err != nil || envVars.BucketName == "" {
		return false, err
	}

	awsRegion := envVars.AWSRegion

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return false, err
	}

	svc := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(envVars.BucketName),
		Key:    aws.String(fileKey),
	}

	_, err2 := svc.DeleteObject(input)

	if err2 != nil {
		return false, err
	}
	return true, nil
}

func (a *AwsInstance) PutImageObjectToS3(file multipart.File, fileHeader *multipart.FileHeader, filePathS3 string) (bool, error) {

	envVars, err := utils.GetEnvVariables()
	if envVars.BucketName == "" {
		return false, err
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
	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal("Error getting file size: ", err)
		return false, err
	}
	fmt.Printf("File size: %d bytes\n", fileSize)

	_, err = file.Seek(0, io.SeekStart)
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

func (a *AwsInstance) RetrieveS3FileInstance() *s3.S3 {
	a.AwsRegion = os.Getenv("AWSREGION")
	a.BucketName = os.Getenv("BUCKETNAME")
	aws_session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(a.BucketName),
		},
	)
	if err != nil {
		a.MyError = err
		return nil
	}
	svc := s3.New(aws_session)
	return svc
}
