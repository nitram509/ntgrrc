package ntgrrc

import (
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const separator = ":"

func storeToken(args *NtgrrcSession, host string, token string) error {
	err := ensureConfigPathExists(args.TokenDir)
	if err != nil {
		return err
	}
	if args.PrintVerbose {
		println("Storing login token " + tokenFilename(args.TokenDir, host))
	}
	data := fmt.Sprintf("%s%s%s", args.model, separator, token)
	return os.WriteFile(tokenFilename(args.TokenDir, host), []byte(data), 0644)
}

func tokenFilename(configDir string, host string) string {
	hash32 := adler32.New()
	io.WriteString(hash32, host)
	return filepath.Join(dotConfigDirName(configDir), "token-"+fmt.Sprintf("%x", hash32.Sum(nil)))
}

func readTokenAndModel2GlobalOptions(args *NtgrrcSession, host string) (NetgearModel, string, error) {
	if len(args.model) > 0 {
		return args.model, args.token, nil
	}
	if args.PrintVerbose {
		println("reading token from: " + tokenFilename(args.TokenDir, host))
	}
	bytes, err := os.ReadFile(tokenFilename(args.TokenDir, host))
	if errors.Is(err, fs.ErrNotExist) {
		return "", "", errors.New("no session (token) exists. please login first")
	}
	data := strings.SplitN(string(bytes), separator, 2)
	if len(data) != 2 {
		return "", "", errors.New("you did an upgrade from a former ntgrcc version. please login again")
	}
	if !isSupportedModel(data[0]) {
		return "", "", errors.New("unknown model stored in token. please login again")
	}
	args.model = NetgearModel(data[0])
	args.token = data[1]
	return args.model, args.token, err
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
