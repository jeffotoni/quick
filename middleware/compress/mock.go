package compress

import (
	"log"
	"net/http"
)

type testcompress struct {
	Request     *http.Request
	HandlerFunc http.HandlerFunc
}

var testcompressSuccess = testcompress{
	Request: &http.Request{
		Header: http.Header{},
	},
	HandlerFunc: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("wr Content-Encoding -> %s", req.Header.Get("Content-Encoding"))
	},
	),
}
