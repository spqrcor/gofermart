package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func FromPostJSON(req *http.Request, input any) error {
	if req.Method != http.MethodPost || !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return errors.New("invalid request")
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		return err
	}
	return nil
}

func FromPostPlain(req *http.Request) (string, error) {
	if req.Method != http.MethodPost || !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
		return "", errors.New("invalid request")
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func IsNumberValid(orderNum string) bool {
	total := 0
	isSecondDigit := false
	for i := len(orderNum) - 1; i >= 0; i-- {
		digit := int(orderNum[i] - '0')
		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		total += digit
		isSecondDigit = !isSecondDigit
	}
	return total%10 == 0
}
