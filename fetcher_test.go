package main

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type fetchMockRet struct {
	r   io.Reader
	err error
}

var _ Fetcher = &fetchMock{}

type fetchMock struct {
	expected map[string]fetchMockRet
}

func (m *fetchMock) Fetch(path string) (io.ReadCloser, error) {
	ret, ok := m.expected[path]
	if !ok {
		panic("unexpected input " + path)
	}
	return ioutil.NopCloser(ret.r), ret.err
}

func TestMultiFetcher(t *testing.T) {
	fetchErr := errors.New("can't fetch")
	msg := "Got it"
	fetcher1 := &fetchMock{
		expected: map[string]fetchMockRet{
			"path": {
				r:   nil,
				err: ErrTryAnother,
			},
			"path2": {
				r:   nil,
				err: fetchErr,
			},
		},
	}
	fetcher2 := &fetchMock{
		expected: map[string]fetchMockRet{
			"path": {
				r:   strings.NewReader(msg),
				err: nil,
			},
			"path2": {
				r:   nil,
				err: ErrTryAnother,
			},
		},
	}
	mf := NewMultiFetcher(fetcher1, fetcher2)
	rc, err := mf.Fetch("path")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer rc.Close()
	res, _ := ioutil.ReadAll(rc)
	strRes := string(res)
	if strRes != msg {
		t.Errorf("bad fetch. Expected: %s, Got: %s", msg, strRes)
	}
	rc, err = mf.Fetch("path2")
	if err != fetchErr {
		t.Errorf("Unexpected error. Expected: %v, Got: %v", fetchErr, err)
	}
}
