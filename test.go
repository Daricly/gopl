package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main()  {
	start := time.Now()
	ch := make(chan string)
	files := os.Args[1:]
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "open file: %v", err)
			os.Exit(1)
		}
		input := bufio.NewScanner(f)
		row := 0
		for input.Scan() {
			row++
			go fetch(input.Text(), ch)
		}
		_ = f.Close()
		for i:=0; i<row; i++ {
			fmt.Println(<-ch)
		}
	}
	//for _, url := range os.Args[1:] {
	//	go fetch(url, ch)
	//}
	//for range os.Args[1:] {
	//	fmt.Println(<-ch)
	//}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
func fetch(url string, ch chan<- string)  {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs	%7d	%s", secs, nbytes, url)
}
