package main

import (
	"log"

	napping "gopkg.in/jmcvetta/napping.v3"
)

// SonyResponse structure.
type SonyResponse interface {
	GetResult() []interface{}
	GetID() int
}

// SonyArrayResponse is...
type SonyArrayResponse struct {
	ID     int           `json:"id"`
	Result []interface{} `json:"result"`
}

// SonyArrayOfArrayResponse is...
type SonyArrayOfArrayResponse struct {
	ID     int             `json:"id"`
	Result [][]interface{} `json:"result"`
}

// GetResult is...
func (r SonyArrayResponse) GetResult() []interface{} {
	return r.Result
}

// GetResult is...
func (r SonyArrayOfArrayResponse) GetResult() []interface{} {
	return r.Result[0]
}

// GetID is...
func (r SonyArrayResponse) GetID() int {
	return r.ID
}

// GetID is...
func (r SonyArrayOfArrayResponse) GetID() int {
	return r.ID
}

// SonyRequest structure.
type SonyRequest struct {
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Version string `json:"version"`
	Params  []int  `json:"params"`
}

func makeRequest(endpoint string, payload *SonyRequest, res SonyResponse) error {
	session := napping.Session{Log: false}
	resp, err := session.Post(APIURL+endpoint, &payload, &res, nil)

	if err != nil {
		log.Fatal(err)
	} else {
		if resp.Status() == 200 {
			return err
		}
	}

	return err
}
