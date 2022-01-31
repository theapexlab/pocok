package aws_utils

import (
	"context"
	"pocok/src/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetAmazonS3Url(bucketName string, region string, key string) string {
	return "https://" + bucketName + ".s3." + region + ".amazonaws.com/" + key
}

func GetAssetUrl(client s3.Client, assetBucketName string, key string) (string, error) {
	region, err := client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: aws.String(assetBucketName),
	})
	if err != nil {
		utils.LogError("Can't find bucket region", err)
		return "", err
	}

	pocokUrl := GetAmazonS3Url(assetBucketName, string(region.LocationConstraint), key)
	return pocokUrl, nil
}
