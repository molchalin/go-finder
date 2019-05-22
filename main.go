package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	pf := ParallelFinder{
		finder: NewFinder(NewMultiFetcher(new(HTTPFetcher), new(FileFetcher))),
	}
	in := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			in <- scanner.Text()
		}
		close(in)
	}()
	out := pf.FindN(1000, in)
	var count int
	for res := range out {
		count += res.n
		fmt.Println(res.String())
	}
	fmt.Printf("Total: %d\n", count)
}
