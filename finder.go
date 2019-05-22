package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type Finder interface {
	FindGo(path string) (int, error)
}

type FinderImpl struct {
	fetcher Fetcher
}

func NewFinder(fetcher Fetcher) Finder {
	return &FinderImpl{
		fetcher: fetcher,
	}
}

func (f *FinderImpl) FindGo(path string) (int, error) {
	rc, err := f.fetcher.Fetch(path)
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

type ParallelFinder struct {
	finder Finder
}

type findRet struct {
	path string
	err  error
	n    int
}

func (r *findRet) String() string {
	msg := fmt.Sprintf("Count for %s: %d", r.path, r.n)
	if r.err != nil {
		msg += fmt.Sprintf(" error: %s", r.err.Error())
	}
	return msg
}

func (f *ParallelFinder) FindN(k uint, in <-chan string) <-chan findRet {
	out := make(chan findRet)
	wg := new(sync.WaitGroup)
	go func() {
		ch := make(chan string)
		var cnt uint
		for v := range in {
			if cnt < k {
				wg.Add(1)
				go func(i uint) {
					defer func() {
						wg.Done()
					}()
					for path := range ch {
						n, err := f.finder.FindGo(path)
						out <- findRet{
							path: v,
							err:  err,
							n:    n,
						}
					}
				}(cnt)
				cnt++
			}
			ch <- v
		}
		close(ch)
		wg.Wait()
		close(out)
	}()
	return out
}
