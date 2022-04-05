package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var wg sync.WaitGroup

func main() {
	//bucket name
	bucket := "ryan-li-practice"
	// Sub floder
	prefix := ""
	// If object Size greater than this value, then compress
	boundarySize := 1024 * 100

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)
	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(prefix)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		if int(*item.Size) > boundarySize {
			wg.Add(1)
			go handler(item)
		}
	}
	wg.Wait()
	fmt.Println("Done")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func handler(item *s3.Object) {

	fmt.Println("Name:         ", *item.Key)
	//	fmt.Println("Last modified:", *item.LastModified)
	//	fmt.Println("Size:         ", *item.Size)
	//	fmt.Println("Storage class:", *item.StorageClass)
	//	fmt.Println("")

	downLoadFile()
	compressFile()
	uploadFile()
	wg.Done()
}

func downLoadFile() {
	fmt.Println("download")
}

func compressFile() {
	fmt.Println("compress")
}

func uploadFile() {
	fmt.Println("Upload")
}
