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

func FetchMock(expected map[string]fetchMockRet) FetchFunc {
	return func(path string) (io.ReadCloser, error) {
		ret, ok := expected[path]
		if !ok {
			panic("unexpected input " + path)
		}
		return ioutil.NopCloser(ret.r), ret.err
	}
}

type findGoTC struct {
	fetch fetchMockRet
	path  string
	count int
	err   error
}

type badReader struct {
	c   int
	err error
}

func (r *badReader) Read(p []byte) (int, error) {
	r.c++
	if r.c == 1 {
		if len(p) > 5 {
			goB := []byte("Go,Go")
			copy(p, goB)
			return 5, nil
		} else {
			panic("wtf")
		}
	}
	return 0, r.err
}

func TestFindGo(t *testing.T) {
	fetchErr := errors.New("can't fetch")
	readErr := errors.New("can't read")
	cases := []findGoTC{
		{
			fetch: fetchMockRet{
				r:   strings.NewReader("Go, Go, Power Rangers!"),
				err: nil,
			},
			path:  "good_path",
			count: 2,
			err:   nil,
		},
		{
			fetch: fetchMockRet{
				r:   nil,
				err: fetchErr,
			},
			path:  "bad_path",
			count: 0,
			err:   fetchErr,
		},
		{
			fetch: fetchMockRet{
				r:   &badReader{0, readErr},
				err: nil,
			},
			path:  "corrupt_path",
			count: 2,
			err:   readErr,
		},
	}

	fetchMap := make(map[string]fetchMockRet)
	for _, c := range cases {
		fetchMap[c.path] = c.fetch
	}
	finder := NewFinder(FetchMock(fetchMap))

	for _, c := range cases {
		n, err := finder.FindGo(c.path)
		if err != c.err {
			t.Errorf("Unexpected err. Expected %v, Got %v", c.err, err)
		}
		if n != c.count {
			t.Errorf("Expected %d Go, got: %d", c.count, n)
		}
	}
}
