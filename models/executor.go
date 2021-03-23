package models

import (
	"net/http"
)

type Req struct {
	Files []string
}

type Response struct {
	File string
}

func (Req) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (i *Req) Bind(r *http.Request) error {
	return nil
}

func (Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (i *Response) Bind(r *http.Request) error {
	return nil
}
