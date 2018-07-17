package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

var GENERATED_KEY_NAME = ".eldon-generated-key"

type SignedMessage struct {
	Message   string
	Signature string
	PubKey    string
}

func (c *SignedMessage) generateSignature(key *ecdsa.PrivateKey, rand io.Reader) error {
	var hash [32]byte = sha256.Sum256([]byte(c.Message))
	bytes, err := key.Sign(rand, hash[:], crypto.SHA256)
	if err != nil {
		return err
	}
	c.Signature = base64.StdEncoding.EncodeToString(bytes)
	marshalled_pubkey, err := x509.MarshalPKIXPublicKey(key.Public())
	c.PubKey = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: marshalled_pubkey}))
	return err
}

type ecdsaSig struct {
	R, S *big.Int
}

func (s *SignedMessage) VerifySignedMessage() error {
	pem, rest := pem.Decode([]byte(s.PubKey))
	if len(rest) != 0 {
		return fmt.Errorf("Too much data in PubKey field")
	}
	if raw_sig, err := base64.StdEncoding.DecodeString(s.Signature); err != nil {
		return err
	} else {
		if key, err := x509.ParsePKIXPublicKey(pem.Bytes); err != nil {
			return err
		} else {
			ecdsaKey, ok := key.(*ecdsa.PublicKey)
			if ok {
				return s.verifyECDSA(ecdsaKey, raw_sig)
			} else {
				return fmt.Errorf("Unknown key type!")
			}
		}
	}
	return nil
}

func (s *SignedMessage) verifyECDSA(key *ecdsa.PublicKey, raw_sig []byte) error {
	var sig ecdsaSig
	rest, err := asn1.Unmarshal(raw_sig, &sig)
	if len(rest) > 0 {
		return fmt.Errorf("Too much data in signature")
	} else if err != nil {
		return err
	}
	var sha_bytes [32]byte = sha256.Sum256([]byte(s.Message))
	if !ecdsa.Verify(key, sha_bytes[:], sig.R, sig.S) {
		return fmt.Errorf("Signature failed to verify")
	}
	return nil
}

func loadKey(pem_bytes []byte) (private_key *ecdsa.PrivateKey, err error) {
	var pem_block *pem.Block
	var rest []byte
	pem_block, rest = pem.Decode(pem_bytes)
	if len(rest) != 0 || pem_block == nil || pem_block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("Improperly generated key file")
	}
	if key, err := x509.ParsePKCS8PrivateKey(pem_block.Bytes); err != nil {
		return nil, err
	} else {
		if key, ok := key.(*ecdsa.PrivateKey); ok {
			return key, nil
		} else {
			return nil, errors.New("Error casting key to ecdsa Private key")
		}
	}
}

func persistKey(private_key *ecdsa.PrivateKey) (private_key_out *ecdsa.PrivateKey, err error) {
	marshalled_key, err := x509.MarshalPKCS8PrivateKey(private_key)
	if err != nil {
		return nil, errors.New("Cannot marshal key to bytes")
	}
	pem_encoded_key := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshalled_key})
	err = ioutil.WriteFile(GENERATED_KEY_NAME, pem_encoded_key, 0600)
	if err != nil {
		return nil, err
	}
	return private_key, nil
}

func loadOrGenerateKey() (private_key *ecdsa.PrivateKey, rand io.ReadCloser, err error) {
	rand, err = os.Open("/dev/urandom")
	if err != nil {
		return nil, nil, errors.New("Cannot open entropy device /dev/urandom")
	}
	var pem_bytes []byte
	if pem_bytes, err = ioutil.ReadFile(GENERATED_KEY_NAME); err == nil {
		key, err := loadKey(pem_bytes)
		return key, rand, err
	} else {
		var new_key *ecdsa.PrivateKey
		new_key, err = ecdsa.GenerateKey(elliptic.P384(), rand)
		if err != nil {
			return nil, nil, errors.New("Cannot generate new key")
		}
		key, err := persistKey(new_key)
		return key, rand, err
	}
}

func main() {
	key, rand, err := loadOrGenerateKey()
	defer rand.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 1 {
		log.Fatal("No message to be signed passed")
	}
	tosign := &SignedMessage{Message: os.Args[1]}
	err = tosign.generateSignature(key, rand)
	if err != nil {
		log.Fatal(err)
	}
	output, err := json.Marshal(tosign)
	if err != nil {
		log.Fatal(err)
	}
	err = tosign.VerifySignedMessage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(output))
}
