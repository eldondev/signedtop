package main

import (
//	"crypto/ecdsa"
//	"crypto/sha256"
	"encoding/json"
	"os"
//	"strings"
	"testing"
)

func TestSignMessage(t *testing.T) {
}

func TestKnownGood(t *testing.T) {
	var file *os.File
	var err error
	if file, err = os.Open("known_good.json"); err != nil {
		t.Errorf("Could not open known good file: %s", err)
	}
	decoder := json.NewDecoder(file)
	var good_message SignedMessage
	if err := decoder.Decode(&good_message); err != nil {
		t.Errorf("Could not read known good file: %s", err)
	}
	err = good_message.VerifySignedMessage();
	if err != nil {
		t.Error(err);
		}
}
