package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func loadToken(args *GlobalOptions) (string, error) {
	if args.Verbose {
		println("reading token from: " + tokenFilename())
	}
	bytes, err := os.ReadFile(tokenFilename())
	return string(bytes), err
}

func requestPage(args *GlobalOptions, url string) (string, error) {
	token, err := loadToken(args)
	if err != nil {
		return "", err
	}

	if args.Verbose {
		println("Fetching data from: " + url)
	}

	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", "SID="+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if args.Verbose {
		println(resp.Status)
	}
	bytes, err := io.ReadAll(resp.Body)
	return string(bytes), err
}
