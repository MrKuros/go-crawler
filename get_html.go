package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	head := res.Header.Get("Content-Type")
	if head == "" {
		return "", fmt.Errorf("missing 'Content-Type' header")
	}

	if !strings.HasPrefix(head, "text/html") {
		return "", fmt.Errorf("expected 'text/html' header, got %s", head)
	}

	statusString := res.Status
	status, _ := strconv.ParseInt(strings.Split(statusString, " ")[0], 10, 64)
	if status > 400 {
		return "", fmt.Errorf("unexpected error, status code %d", status)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(body), nil
}
