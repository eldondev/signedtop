package main

import (
	"context"
	"errors"
	"github.com/eldondev/signedtop"
	"net/http"
)

type SignedTopServer struct{ blocks []signedtop.Block }

func (s *SignedTopServer) GetBlock(context.Context, *signedtop.BlockId) (*signedtop.Block, error) {
	return nil, errors.New("Not implemented")
}

func (s *SignedTopServer) GetTop(context.Context, *signedtop.EmptyParams) (*signedtop.Block, error) {
	return nil, errors.New("Not implemented")
}

func (s *SignedTopServer) GetPubKey(context.Context, *signedtop.EmptyParams) (*signedtop.PubKey, error) {
	return nil, errors.New("Not implemented")
}

func (s *SignedTopServer) PleaseSign(context.Context, *signedtop.DataToSign) (*signedtop.Block, error) {
	return nil, errors.New("Not implemented")
}

// Run the implementation in a local server
func main() {
	signerHandler := signedtop.NewSignedTopServer(&SignedTopServer{}, nil)
	// You can use any mux you like - NewHelloWorldServer gives you an http.Handler.
	mux := http.NewServeMux()
	// The generated code includes a const, <ServiceName>PathPrefix, which
	// can be used to mount your service on a mux.
	mux.Handle(signedtop.SignedTopPathPrefix, signerHandler)
	http.ListenAndServe(":8080", mux)
}
