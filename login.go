package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"math"
	"strings"
)

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
