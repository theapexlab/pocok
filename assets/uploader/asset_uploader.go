package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

const (
	ASSET_FOLDER  = "./assets/"
	AWS_S3_REGION = "eu-central-1"
)

func main() {

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(ASSET_FOLDER)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start uploading assets.")
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			log.Println("uploading - " + fileName)
			err = uploadFile(session, fileName)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Uploading assets completed.")
}

func uploadFile(session *session.Session, fileName string) error {

	upFile, err := os.Open(ASSET_FOLDER + fileName)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_ASSET_BUCKET_NAME")),
		Key:           aws.String(fileName),
		ACL:           aws.String("public-read"),
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
	})
	return err
}
