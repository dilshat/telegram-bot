package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/labstack/gommon/log"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func FileExists(name string) bool {
	_, err := os.Stat(name)

	if os.IsNotExist(err) {
		return false
	}

	//sometimes there can be permission or other errors
	//here we use a simple logic that if file exists and we can use it then true otherwise false
	return err == nil
}

type FileType int

const (
	PHOTO FileType = iota
	AUDIO
	VIDEO
	OTHER
)

func GetFileType(path string) FileType {
	file, err := os.Open(path)
	if err != nil {
		log.Error("Failed to determine file type of "+path, err)
		return OTHER
	}
	defer file.Close()

	head := make([]byte, 261)
	if _, err = file.Read(head); err != nil {
		log.Error("Failed to read file header of "+path, err)
		return OTHER
	}

	if filetype.IsImage(head) {
		return PHOTO
	} else if filetype.IsVideo(head) {
		return VIDEO
	} else if filetype.IsAudio(head) {
		return AUDIO
	}

	return OTHER
}

func ParseFileType(fileType string) FileType {
	fileType = strings.TrimSpace(strings.ToLower(fileType))
	switch fileType {
	case "video":
		return VIDEO
	case "audio":
		return AUDIO
	case "photo":
		return PHOTO
	default:
		return OTHER
	}
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func doPOST(aURL string, payload map[string]interface{}, headers map[string]interface{}) (string, error) {
	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest("POST", aURL, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	for key, val := range headers {
		req.Header.Set(key, fmt.Sprintf("%v", val))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func doGET(aURL string, params map[string]interface{}, headers map[string]interface{}) (string, error) {
	urlParams := url.Values{}

	for key, val := range params {
		urlParams.Add(key, fmt.Sprintf("%v", val))
	}

	req, err := http.NewRequest("GET", aURL, nil)
	if err != nil {
		return "", err
	}
	for key, val := range headers {
		req.Header.Set(key, fmt.Sprintf("%v", val))
	}

	req.URL.RawQuery = urlParams.Encode()

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ReadFile(path string) (string, error) {
	if !FileExists(path) {
		return "", errors.New(fmt.Sprintf("File %s not found", path))
	}

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}
