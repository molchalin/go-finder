package main

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
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
	finder := NewFinder(&fetchMock{fetchMap})

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

type finderMockRet struct {
	n   int
	err error
}

var _ Finder = &finderMock{}

type finderMock struct {
	expected map[string]finderMockRet
}

func (m *finderMock) FindGo(path string) (int, error) {
	ret, ok := m.expected[path]
	if !ok {
		panic("unexpected input " + path)
	}
	return ret.n, ret.err
}

func TestParallel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	go func() {
		defer cancel()
		templateOut := findRet{
			path: "path",
			err:  nil,
			n:    10,
		}
		k := uint(1)
		in := make(chan string)
		finder := &finderMock{
			expected: map[string]finderMockRet{
				templateOut.path: {
					n:   templateOut.n,
					err: templateOut.err,
				},
			},
		}
		pf := ParallelFinder{
			finder: finder,
		}
		out := pf.FindN(k, in)
		go func() {
			for i := 0; i < 5; i++ {
				in <- "path"
			}
			close(in)
		}()
		for res := range out {
			if !reflect.DeepEqual(res, templateOut) {
				t.Errorf("Expected %#v, Got %#v", templateOut, res)
			}
		}
	}()
	<-ctx.Done()
	if err := ctx.Err(); err == context.DeadlineExceeded {
		t.Errorf(err.Error())
	}
}
