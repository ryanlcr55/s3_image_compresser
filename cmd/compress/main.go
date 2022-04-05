package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var wg sync.WaitGroup
var region string = "ap-northeast-1"
var bucket string = "ryan-li-practice"

// Sub floder
var prefix string = ""

// If object Size greater than this value, then compress
var boundarySize int = 1024 * 100

func main() {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

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
		if int(*item.Size) > boundarySize {
			wg.Add(1)
			go handler(sess, item)
		}
	}
	wg.Wait()
	fmt.Println("Done")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func handler(sess *session.Session, item *s3.Object) {
	defer wg.Done()
	fmt.Println("Name:         ", *item.Key)
	//	fmt.Println("Last modified:", *item.LastModified)
	//	fmt.Println("Size:         ", *item.Size)
	//	fmt.Println("Storage class:", *item.StorageClass)
	//	fmt.Println("")

	downLoadFile(sess, *item.Key)
}

func downLoadFile(sess *session.Session, key string) {
	fmt.Println("download")
	svc := s3.New(sess)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	out, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		exitErrorf("Unable to Down File Key: %q, %v", key, err)
	}

	defer out.Body.Close()
}

func compressFile() {
	fmt.Println("compress")
}

func uploadFile() {
	fmt.Println("Upload")
}
