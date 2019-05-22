package main

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type mockRet struct {
	str string
	err error
}

func FetchMock(expected map[string]mockRet) FetchFunc {
	return func(path string) (io.ReadCloser, error) {
		ret, ok := expected[path]
		if !ok {
			panic("unexpected input " + path)
		}
		return ioutil.NopCloser(strings.NewReader(ret.str)), ret.err
	}
}

func TestFindGo(t *testing.T) {
	path := "nowhere"
	fetcher := FetchMock(map[string]mockRet{
		path: mockRet{
			str: "Go, Go, Power Rangers!",
		},
	})
	finder := NewFinder(fetcher)
	n, err := finder.FindGo(path)
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}
	if n != 2 {
		t.Errorf("Expected 2 go, got: %d", n)
	}
}
