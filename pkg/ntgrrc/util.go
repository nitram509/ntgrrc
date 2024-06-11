package ntgrrc

import (
	"errors"
	"strconv"
	"strings"
)

func suffixToLength(s string, length int) string {
	if len(s) < length {
		diff := length - len(s)
		return s + strings.Repeat(" ", diff)
	}
	return s
}

func parseFloat32(text string) float32 {
	i64, _ := strconv.ParseFloat(text, 32)
	return float32(i64)
}

func parseInt32(text string) int32 {
	i64, _ := strconv.ParseInt(text, 10, 32)
	return int32(i64)
}

func ensureModelIs30x(args *NtgrrcSession, host string) error {
	model, _, err := readTokenAndModel2GlobalOptions(args, host)
	if err != nil {
		return err
	}
	if !isModel30x(model) {
		return errors.New("This command is not yet supported for your Netgear model. " +
			"You might want to support the project by creating an issue on Github")
	}
	return nil
}
