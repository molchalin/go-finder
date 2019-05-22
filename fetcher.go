package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

var ErrTryAnother = errors.New("this fetcher can't handle this path")

type Fetcher interface {
	Fetch(string) (io.ReadCloser, error)
}

type MultiFetcher struct {
	fetchers []Fetcher
}

func NewMultiFetcher(fetchers ...Fetcher) *MultiFetcher {
	return &MultiFetcher{
		fetchers: fetchers,
	}
}

func (mf *MultiFetcher) Fetch(path string) (io.ReadCloser, error) {
	for _, f := range mf.fetchers {
		rc, err := f.Fetch(path)
		if err != nil {
			if err != ErrTryAnother {
				return nil, err
			}
		} else {
			return rc, nil
		}
	}
	return nil, ErrTryAnother
}

type FileFetcher struct{}

func (ff *FileFetcher) Fetch(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

type HTTPFetcher struct{}

func (hf *HTTPFetcher) Fetch(path string) (io.ReadCloser, error) {
	if _, err := url.Parse(path); err != nil {
		return nil, ErrTryAnother
	}
	req, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	return req.Body, nil
}
