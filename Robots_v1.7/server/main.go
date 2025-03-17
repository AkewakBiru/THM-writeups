package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		fmt.Printf("[%s] %s %s\n", clientIP, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func test_sanity(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	fmt.Printf("[%s] %s %s\n", clientIP, r.Method, r.URL.Path)

	w.Write([]byte("test response"))
}

func redir_me(w http.ResponseWriter, r *http.Request) {
	reqt, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(reqt))

	w.Header().Set("Location", "http://localhost:9000")
	w.WriteHeader(http.StatusMovedPermanently)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	if r.Method == http.MethodPost {
		rd, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		plain, err := url.QueryUnescape(string(rd))
		if err != nil {
			log.Fatal(err)
		}

		if start := strings.Index(plain, "<body>"); start != -1 {
			fmt.Printf("Read: %s\n", plain[start+len("<body>"):])
		} else {
			fmt.Printf("Body tag not found\n")
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	fs := http.StripPrefix("/f/", http.FileServer(http.Dir("/Users/akewakbiru/Documents/try_hack_me_git/THM-writeups/Robots_v1.7")))
	http.Handle("/f/", LoggingMiddleware(fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/redir_me", redir_me)
	http.HandleFunc("/test", test_sanity)
	log.Printf("Server started listening at 0.0.0.0:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
