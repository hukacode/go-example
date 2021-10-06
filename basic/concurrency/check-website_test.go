package concurrency

import (
	"reflect"
	"testing"
)

func mockWebsiteChecker(url string) bool {
	if url == "abc://test.com" {
		return false
	}

	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
	}

	want := map[string]bool{
		"http://google.com": true,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}
