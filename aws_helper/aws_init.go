package awshelper

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

type AwsInstance struct {
	BucketName string 
	AwsRegion string
	MyError error
}
 


func (a *AwsInstance) AwsInit()(*session.Session, error){
	err :=	godotenv.Load(".env")
	if err !=nil { 
		a.MyError = err
		return nil, a.MyError
	}

	a.AwsRegion = os.Getenv("AWSREGION") 
	a.BucketName = os.Getenv("BUCKETNAME") 

	aws_sesson, err := session.NewSession(
		&aws.Config{
			Region: aws.String(a.BucketName),
		},
	)

	if err !=nil{ 
		a.MyError = err
		return nil, a.MyError
	}

	return aws_sesson, a.MyError
}

func (a *AwsInstance) PutImageObjectToS3(file multipart.File, fileHeader *multipart.FileHeader,filePath string)(bool, error){
	err := godotenv.Load(".env")
	if err !=nil{
		log.Fatal("Error in loading Environments",err)
		return false,err
	}
	awsRegion := os.Getenv("AWSREGION") 
	awsBucket := os.Getenv("BUCKETNAME") 

	sess, err :=  session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil{
		log.Fatal("Error creating session: ",err)
		return false ,err
	}

	_,err = io.ReadAll(file)
	if err != nil {
		log.Fatal("Error in reading file: ",err)
		return false ,err
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

 
	input :=  s3.PutObjectInput{
		Bucket:      aws.String(awsBucket),
        Key:         aws.String(path.Join("wallpapers/",fileHeader.Filename)),
        Body:        file,
        ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	}

	_ ,err = svc.PutObject(&input) 
	if err !=nil{
		log.Fatal("Error in Uploading file: " ,err)
		return false ,err
	}

	fmt.Printf("File Uploaded Successfully!")
	return true ,err
}


func (a *AwsInstance) RetrieveS3FileInstance()(*s3.S3){
	a.AwsRegion = os.Getenv("AWSREGION") 
	a.BucketName = os.Getenv("BUCKETNAME") 
	aws_session, err :=  session.NewSession(
		&aws.Config{
			Region: aws.String(a.BucketName),
		},
	)
	if err !=nil{
		a.MyError = err
		return nil
	}
	svc := s3.New(aws_session)
	return svc
}