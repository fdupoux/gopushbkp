package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
	"strings"
)

func uploadArchive(archfile string) error {
	fmt.Printf("Uploading archive %s to AWS S3 ...\n", archfile)

	// Pass credentials via environment variables
	os.Setenv("AWS_ACCESS_KEY_ID", config.aws_access)
	os.Setenv("AWS_SECRET_ACCESS_KEY", config.aws_secret)

	// Create a single AWS session (we can re use this if we're uploading many files)
	mysession, err := session.NewSession(&aws.Config{Region: aws.String(config.aws_region)})
	if err != nil {
		return fmt.Errorf("Failed to create AWS Session: %s", err)
	}

	// Upload
	files := []string{archfile, archfile + ".sha256"}
	for _, curfile := range files {
		err = uploadFile(mysession, curfile, config.aws_prefix)
		if err != nil {
			return fmt.Errorf("Failed to upload file %s to S3: %s", curfile, err)
		}
	}

	return nil
}

func uploadFile(mysession *session.Session, fullpath string, prefix string) error {

	// Open the file for use
	myfile, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %v", fullpath, err)
	}
	defer myfile.Close()

	// Get file size and read the file content into a buffer
	basename := filepath.Base(fullpath)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(mysession)

	// Determine name of the key in S3
	nameparts := []string{strings.Trim(prefix, "/"), basename}
	objname := strings.Join(nameparts, "/")

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.aws_bucket),
		Key:    aws.String(objname),
		Body:   myfile,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Upload to %v completed\n", result.Location)

	return nil
}
