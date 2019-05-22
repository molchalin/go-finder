package main

import (
	"bufio"
	"io"
	"strings"
)

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
	rc, err := f.fetch(path)
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	scanner := bufio.NewScanner(rc)
	scanner.Split(bufio.ScanWords)

	var count int
	for scanner.Scan() {
		count += strings.Count(scanner.Text(), "Go")
	}
	return count, scanner.Err()
}

type findRet struct {
	path  string
	err   error
	count int
}

func (f *Finder) findN(k uint, in <-chan string, out chan<- findRet) {
	return
}
