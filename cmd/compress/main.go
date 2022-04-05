package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/h2non/bimg"
)

const (
	region string = "ap-northeast-1"
	bucket string = "ryan-li-practice"
	// Sub floder
	prefix string = ""

	// If object Size greater than this value, then compress
	boundarySize int = 1024 * 100
)

var wg sync.WaitGroup

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
	fmt.Println("Processing Name:         ", *item.Key)

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
	buffer, bufferErr := io.ReadAll(out.Body)
	if bufferErr != nil {
		exitErrorf("Cant buffer the file", key, err)
	}

	processed := compressFile(buffer)
	fmt.Println(processed)

	defer out.Body.Close()
}

func compressFile(buffer []byte) []byte {
	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		exitErrorf("1111")
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 75})
	if err != nil {
		exitErrorf("2222")
	}

	return processed
}
