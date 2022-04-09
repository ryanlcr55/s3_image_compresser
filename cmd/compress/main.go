package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/h2non/bimg"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	region := os.Getenv("REGION")
	bucket := os.Getenv("BUCKET_NAME")
	prefix := os.Getenv("PREFIX")
	boundarySize := os.Getenv("BOUNDARY_SIZE")
	boundarySizeKB, _ := strconv.Atoi(boundarySize)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	// Create S3 service client
	svc := s3.New(sess)
	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		if int(*item.Size) > boundarySizeKB {
			handler(sess, item)
		}
	}
	fmt.Println("Done")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func handler(sess *session.Session, item *s3.Object) {
	fmt.Println("Processing Name:         ", *item.Key)

	key := *item.Key
	svc := s3.New(sess)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	// Get file from S3
	out, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(key),
	})
	if err != nil {
		exitErrorf("Unable to Download File Key: %q, %v", key, err)
	}
	buffer, bufferErr := io.ReadAll(out.Body)
	if bufferErr != nil {
		exitErrorf("Unable to buffer the file", key, err)
	}

	compressFile(&buffer)

	// Update the compressed File
	_, putFileErr := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer),
	})
	if putFileErr != nil {
		exitErrorf(" buffer fail", key, err)
	}
	defer out.Body.Close()
}

func compressFile(buffer *[]byte) {
	converted, err := bimg.NewImage(*buffer).Convert(bimg.WEBP)
	if err != nil {
		exitErrorf("Covert Failed")
	}

	*buffer, err = bimg.NewImage(converted).Process(bimg.Options{Quality: 75})
	if err != nil {
		exitErrorf("Process Failed")
	}
}
