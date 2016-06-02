package main

import (
	"testing"
	"time"
)

func TestEndToEndProcess(t *testing.T) {
	run(
		time.Now(),
		"http://localhost:49153",
		"localhost:1080",
		"TEST",
		"TEST",
		"/data/source",
		"/data/target",
		"domain",
		"bucketName",
		"env",
	)
	t.Fail()
}