package main

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"pocok/src/utils/aws_clients"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	client := aws_clients.GetS3Client()

	_, filename, _, _ := runtime.Caller(0)
	assetFolder := path.Dir(filename) + "/../../assets/"

	assetBucketName, err := getAssetBucketName(client)
	if err != nil {
		log.Fatal(err)
	}

	emptyAssetBucket(client, assetBucketName)

	files, err := ioutil.ReadDir(assetFolder)
	if err != nil {
		log.Fatal(err)
	}
	uploadFiles(client, assetBucketName, assetFolder, files)

}

func uploadFiles(client *s3.Client, assetBucketName string, assetFolder string, files []fs.FileInfo) {
	log.Println("Start uploading assets.")
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			log.Println("uploading - " + fileName)
			err := uploadObject(client, assetBucketName, assetFolder, fileName)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Println("Uploading assets completed.")
}

func getAssetBucketName(client *s3.Client) (string, error) {
	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	assetBucketName := ""
	for _, bucket := range buckets.Buckets {
		isAssetBucket := strings.Contains(*bucket.Name, "assetbucket")
		if isAssetBucket {
			assetBucketName = *bucket.Name
		}
	}

	if assetBucketName == "" {
		return "", errors.New("can't find asset bucket")
	}

	return assetBucketName, nil
}

func uploadObject(client *s3.Client, assetBucketName string, assetFolder string, fileName string) error {
	upFile, err := os.Open(assetFolder + fileName)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(assetBucketName),
		Key:           aws.String(fileName),
		ACL:           "public-read",
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: fileSize,
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
	})

	return err
}

func deleteObject(client *s3.Client, bucket, key, versionId *string) {
	log.Println("deleting - " + *key)
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket:    bucket,
		Key:       key,
		VersionId: versionId,
	})
	if err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}
}

func emptyAssetBucket(client *s3.Client, assetBucketName string) {
	log.Println("Start deleting items in asset bucket.")
	bucket := aws.String(assetBucketName)

	in := &s3.ListObjectsV2Input{Bucket: bucket}
	for {
		out, err := client.ListObjectsV2(context.TODO(), in)
		if err != nil {
			log.Fatalf("Failed to list objects: %v", err)
		}

		for _, item := range out.Contents {
			deleteObject(client, bucket, item.Key, nil)
		}

		if out.IsTruncated {
			in.ContinuationToken = out.ContinuationToken
		} else {
			break
		}
	}

	inVer := &s3.ListObjectVersionsInput{Bucket: bucket}
	for {
		out, err := client.ListObjectVersions(context.TODO(), inVer)
		if err != nil {
			log.Fatalf("Failed to list version objects: %v", err)
		}

		for _, item := range out.DeleteMarkers {
			deleteObject(client, bucket, item.Key, item.VersionId)
		}

		for _, item := range out.Versions {
			deleteObject(client, bucket, item.Key, item.VersionId)
		}

		if out.IsTruncated {
			inVer.VersionIdMarker = out.NextVersionIdMarker
			inVer.KeyMarker = out.NextKeyMarker
		} else {
			break
		}
	}
	log.Println("Deleting assets in asset bucket completed.")
}
