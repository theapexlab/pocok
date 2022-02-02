package main

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"pocok/src/utils"
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

	assetBucketName, getBucketError := getAssetBucketName(client)
	if getBucketError != nil {
		utils.LogFatal(getBucketError)
	}

	emptyAssetBucket(client, assetBucketName)

	files, readDirError := ioutil.ReadDir(assetFolder)
	if readDirError != nil {
		utils.LogFatal(readDirError)
	}
	uploadFiles(client, assetBucketName, assetFolder, files)

}

func uploadFiles(client *s3.Client, assetBucketName string, assetFolder string, files []fs.FileInfo) {
	utils.Log("Start uploading assets.")
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			utils.Logf("Uploading - %s \n", fileName)
			uploadError := uploadObject(client, assetBucketName, assetFolder, fileName)
			if uploadError != nil {
				utils.LogFatal(uploadError)
			}
		}
	}
	utils.Log("Uploading assets completed.")
}

func getAssetBucketName(client *s3.Client) (string, error) {
	buckets, listBucketError := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if listBucketError != nil {
		utils.LogFatal(listBucketError)
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
	upFile, fileOpenError := os.Open(assetFolder + fileName)
	if fileOpenError != nil {
		return fileOpenError
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	_, readError := upFile.Read(fileBuffer)
	if readError != nil {
		return readError
	}

	_, s3UploadError := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(assetBucketName),
		Key:           aws.String(fileName),
		ACL:           "public-read",
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: fileSize,
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
	})

	return s3UploadError
}

func deleteObject(client *s3.Client, bucket, key, versionId *string) {
	utils.Logf("Deleting - %s \n", *key)
	_, s3DeleteError := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket:    bucket,
		Key:       key,
		VersionId: versionId,
	})
	if s3DeleteError != nil {
		utils.LogFatalf("Failed to delete object: %v", s3DeleteError)
	}
}

func emptyAssetBucket(client *s3.Client, assetBucketName string) {
	utils.Log("Start deleting items in asset bucket.")
	bucket := aws.String(assetBucketName)

	in := &s3.ListObjectsV2Input{Bucket: bucket}
	for {
		out, s3ListError := client.ListObjectsV2(context.TODO(), in)
		if s3ListError != nil {
			utils.LogFatalf("Failed to list objects: %v", s3ListError)
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
		out, s3ListObjectVersionsError := client.ListObjectVersions(context.TODO(), inVer)
		if s3ListObjectVersionsError != nil {
			utils.LogFatalf("Failed to list version objects: %v", s3ListObjectVersionsError)
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
	utils.Log("Deleting assets in asset bucket completed.")
}
