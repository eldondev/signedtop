package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
)

type SignedMessage struct {
	Message   string
	Signature string
	PubKey    string
}

func (c *SignedMessage) generateSignature() error {
	return errors.New("Not yet implemented!")
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

func main() {

	log.Printf("%+v", (&SignedMessage{Message: os.Args[1]}).generateSignature())
}
