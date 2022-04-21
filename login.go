package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const FailedAttempt = "no SID cookie found in response header"

type LoginCommand struct {
	Address  string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Password string `required:"" help:"the admin console's password'" short:"p"`
}

func (login *LoginCommand) Run(args *GlobalOptions) error {
	seedValue, err := getSeedValueFromSwitch(args, login.Address)
	if err != nil {
		return err
	}

	encryptedPwd := encryptPassword(login.Password, seedValue)

	err = doLogin(args, login.Address, encryptedPwd)
	if err != nil {
		return err
	}

	return nil
}

func doLogin(args *GlobalOptions, host string, encryptedPwd string) error {
	url := fmt.Sprintf("http://%s/login.cgi", host)
	if args.Verbose {
		println("login attempt: " + url)
	}
	formData := "password=" + encryptedPwd
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(formData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if args.Verbose {
		println(resp.Status)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	token := getSessionToken(resp)
	if token == FailedAttempt && resp.StatusCode == http.StatusOK {
		return errors.New("login request returned 200 OK, but response did not contain a session token (SID cookie value#). " +
			"this is known behaviour from the switch. please, wait some minutes and tray again later")
	}

	err = ensureConfigPathExists()
	if err != nil {
		return err
	}
	if args.Verbose {
		println("Storing login token " + tokenFilename())
	}
	err = os.WriteFile(tokenFilename(), []byte(token), 0644)
	if err != nil {
		return err
	}
	return nil
}

func tokenFilename() string {
	return filepath.Join(dotConfigDirName(), "token")
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

func getSessionToken(resp *http.Response) string {
	cookie := resp.Header.Get("Set-Cookie")
	const SessionIdPrefix = "SID="
	if strings.HasPrefix(cookie, SessionIdPrefix) {
		sidVal := cookie[len(SessionIdPrefix):]
		split := strings.Split(sidVal, ";")
		return split[0]
	}
	return FailedAttempt
}

func getSeedValueFromSwitch(args *GlobalOptions, host string) (string, error) {
	url := fmt.Sprintf("http://%s/login.cgi", host)
	if args.Verbose {
		println("fetch seed value from: " + url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if args.Verbose {
		println(resp.Status)
	}
	defer resp.Body.Close()

	seedValue, err := getSeedValueFromLoginHtml(resp.Body)
	if err != nil {
		return "", err
	}
	return seedValue, nil
}

func getSeedValueFromLoginHtml(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}
	randVal, exists := doc.Find("#rand").First().Attr("value")

	if exists {
		return randVal, nil
	}
	return "", errors.New("random seed value not found in login.cgi response. " +
		"An element with id=rand and an attribute 'value' is expected")
}

// encryptPassword re-implements some logic from Netgear's GS305EP frontend component, see login.js
func encryptPassword(password string, seedValue string) string {
	mergedStr := specialMerge(password, seedValue)
	hash := md5.New()
	io.WriteString(hash, mergedStr)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func specialMerge(password string, seedValue string) string {
	result := strings.Builder{}
	maxLen := int(math.Max(float64(len(password)), float64(len(seedValue))))
	for i := 0; i < maxLen; i++ {
		if i < len(password) {
			result.WriteString(string([]rune(password)[i]))
		}
		if i < len(seedValue) {
			result.WriteString(string([]rune(seedValue)[i]))
		}
	}
	return result.String()
}
