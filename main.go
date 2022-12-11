package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	list       []string = ReadAllLines("ips.txt")
	file       *os.File
	checked    int
	found      int
	rs, realrs int
	fiuler     = make(chan string)
)

func main() {
	fmt.Printf("Loaded IP's => %d\n", len(list))
	o := make(chan string)
	var err error
	file, err = os.Create("result.txt")
	if err != nil {
		panic(err)
	}

	t := 0
	fmt.Print("Threads: ")
	fmt.Scanln(&t)
	go func() {
		for {
			result := <-fiuler
			fmt.Print("\033[2K" + result)
			file.WriteString(result)
		}
	}()
	go func() {
		for {
			realrs = rs
			rs = 0
			time.Sleep(time.Second)
		}
	}()
	for i := 0; i != t; i++ {
		go thread(o)
	}
	go func() {
		for {
			fmt.Printf("\rChecked: %d, C/s: %d, Found: %d        ", checked, realrs, found)
			time.Sleep(time.Millisecond * 188)
		}
	}()
	for _, line := range list {
		o <- line
	}
	for {
		time.Sleep(time.Millisecond)
	}
}

func thread(o chan string) {
	for {
		ip := <-o
		domains := gedDmn(ip)
		checked++
		rs++
		for _, domain := range domains {
			found++
			fiuler <- (fmt.Sprintf("[%s] => [%s]\r\n", ip, domain))
		}
	}
}

func gedDmn(ip string) []string {
	for tries := 0; tries < 5; tries++ {
		dmns, err := net.LookupAddr(ip)
		if err == nil {
			return dmns
		}
	}
	return nil
}

func ReadAllLines(path string) (lines []string) {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
