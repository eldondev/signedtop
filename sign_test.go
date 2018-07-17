package main

import (
	"bytes"
	//	"crypto/ecdsa"
	//	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
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
		t.Errorf("Could not open known_good file: %s", err)
	}
	decoder := json.NewDecoder(file)
	var good_message SignedMessage
	if err := decoder.Decode(&good_message); err != nil {
		t.Errorf("Could not decode json from known_good file: %s", err)
	}
	err = good_message.VerifySignedMessage()
	if err != nil {
		t.Error(err)
	}
}

func TestMainRaw(t *testing.T) {
	os.Args = []string{"bump", "thump"};
	main()
}

func TestMainPerservePkey(t *testing.T) {
	tmp, tmp_err := ioutil.TempFile("", "key")
	os.Args = []string{"bump", "thump"};
	os.Remove(tmp.Name())
	GENERATED_KEY_NAME = tmp.Name()
	main()
	key, err := ioutil.ReadFile(GENERATED_KEY_NAME)
	if len(key) == 0 || err != nil || tmp_err != nil {
		t.Fatalf("No key generated")
	}
	main()
	key2, err := ioutil.ReadFile(GENERATED_KEY_NAME)
	if len(key2) == 0 || err != nil || tmp_err != nil {
		t.Fatalf("No key generated")
	}
	if bytes.Compare(key, key2) != 0 {
		t.Fatalf("Key not reused:\n %+v,\n %+v", key, key2)
	}
}
