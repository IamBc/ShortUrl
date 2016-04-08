package main

import (
	"testing"
)

func TestInitStorage(t *testing.T) {
	InitStorage()
}

func TestGetURLFromStorage(t *testing.T) {
	InitStorage()

	url, err := GetURLFromStorage(`jH1ufr87`)

	if err != nil || url != `http://www.dir.bg` {
		t.Fail()
	}
}
