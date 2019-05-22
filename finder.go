package main

import "io"

type FetchFunc func(string) (io.ReadCloser, error)
type Finder struct {
	fetch FetchFunc
}

func NewFinder(fetch FetchFunc) *Finder {
	return &Finder{
		fetch: fetch,
	}
}

func (f *Finder) FindGo(path string) (int, error) {
	return 0, nil
}

func (f *Finder) FindAllGo(r io.Reader, w io.Writer) (int, error) {
	return 0, nil
}
