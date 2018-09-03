package utils

import (
	r "github.com/levigross/grequests"
)

type RequestData r.RequestOptions

func HttpGet(url string, options ...RequestData) (*r.Response, error) {
	if len(options) > 0 {
		data := r.RequestOptions(options[0])
		return r.Get(url, &data)
	}
	return r.Get(url, nil)
}

func HttpPost(url string, options RequestData) (*r.Response, error) {
	data := r.RequestOptions(options)
	return r.Post(url, &data)
}
