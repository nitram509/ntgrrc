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
	err := ensureConfigPathExists(args.TokenDir)
	if err != nil {
		return err
	}
	if args.Verbose {
		println("Storing login token " + tokenFilename(args.TokenDir, host))
	}
	return os.WriteFile(tokenFilename(args.TokenDir, host), []byte(token), 0644)
}

func tokenFilename(configDir string, host string) string {
	hash32 := adler32.New()
	io.WriteString(hash32, host)
	return filepath.Join(dotConfigDirName(configDir), "token-"+fmt.Sprintf("%x", hash32.Sum(nil)))
}

func loadToken(args *GlobalOptions, host string) (string, error) {
	if args.Verbose {
		println("reading token from: " + tokenFilename(args.TokenDir, host))
	}
	bytes, err := os.ReadFile(tokenFilename(args.TokenDir, host))
	if errors.Is(err, fs.ErrNotExist) {
		return "", errors.New("no session (token) exists. please login first")
	}
	return string(bytes), err
}

func ensureConfigPathExists(configDir string) error {
	dotConfigNtgrrc := dotConfigDirName(configDir)
	err := os.MkdirAll(dotConfigNtgrrc, os.ModeDir|0700)
	return err
}

func dotConfigDirName(configDir string) string {
	if configDir == "" {
		configDir = os.TempDir()
	}
	return filepath.Join(configDir, ".config", "ntgrrc")
}
