package main

import (
	"testing"
)

func TestEndToEndProcess(t *testing.T) {
	mockFleet := mockFleetAPI{}
	runInner(mockFleet)
}
