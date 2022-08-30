package main

import (
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func storeToken(args *GlobalOptions, host string, token string) error {
	err := ensureConfigPathExists()
	if err != nil {
		return err
	}
	if args.Verbose {
		println("Storing login token " + tokenFilename(host))
	}
	return os.WriteFile(tokenFilename(host), []byte(token), 0644)
}

func tokenFilename(host string) string {
	hash32 := adler32.New()
	io.WriteString(hash32, host)
	return filepath.Join(dotConfigDirName(), "token-"+fmt.Sprintf("%x", hash32.Sum(nil)))
}

func loadToken(args *GlobalOptions, host string) (string, error) {
	if args.Verbose {
		println("reading token from: " + tokenFilename(host))
	}
	bytes, err := os.ReadFile(tokenFilename(host))
	if errors.Is(err, fs.ErrNotExist) {
		return "", errors.New("no session (token) exists. please login first")
	}
	return string(bytes), err
}

func ensureConfigPathExists() error {
	dotConfigNtgrrc := dotConfigDirName()
	err := os.MkdirAll(dotConfigNtgrrc, os.ModeDir)
	return err
}

func dotConfigDirName() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".config", "ntgrrc")
}
