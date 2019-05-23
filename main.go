package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	k := flag.Uint("k", 3, "parallel factor")
	flag.Parse()
	if *k <= 0 {
		panic("k must be greater than 0")
	}
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
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		close(in)
	}()
	out := pf.FindN(*k, in)
	var count int
	for res := range out {
		count += res.n
		fmt.Println(res.String())
	}
	fmt.Printf("Total: %d\n", count)
}
