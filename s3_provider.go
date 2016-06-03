package main

import (
	"io"
	"github.com/rlmcpherson/s3gof3r"
)

// S3WriterProvider AWS S3 writer provider
type S3WriterProvider struct {
	bucket *s3gof3r.Bucket
}

func newS3WriterProvider(awsAccessKey string, awsSecretKey string, s3Domain string, bucketName string) *S3WriterProvider {
	s3gof3r.DefaultDomain = s3Domain

	awsKeys := s3gof3r.Keys{
		AccessKey: awsAccessKey,
		SecretKey: awsSecretKey,
	}

	s3 := s3gof3r.New("", awsKeys)
	bucket := s3.Bucket(bucketName)

	return &S3WriterProvider{bucket}
}

func (writerProvider *S3WriterProvider) getWriter(fileName string) (io.WriteCloser, error) {
	return writerProvider.bucket.PutWriter(fileName, nil, nil)
}
