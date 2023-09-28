package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// src: https://golangcode.com/generate-sha256-hmac/
func hmacsha256(data string, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(data))
	sha := hex.EncodeToString(hash.Sum(nil))
	return sha
}

func getApiEndpoint(api string) string {
	var endpoint string

	apiRegex := regexp.MustCompile(`(/api/.*)`)
	versionRegex := regexp.MustCompile(`(^v1/.*)`)

	var apiMatch []string
	apiMatch = apiRegex.FindStringSubmatch(endpoint)

	if len(apiMatch) == 0 {
		apiMatch = versionRegex.FindStringSubmatch(endpoint)
	}
	if len(apiMatch) > 1 {
		endpoint = apiMatch[1]
	}

	endpoint = regexp.MustCompile(`[?&][^=?&]+=[^&]+`).ReplaceAllString(endpoint, "")

	return endpoint
}

func getTime(tz string) (string, string) {
	location, err := time.LoadLocation(tz)

	if err != nil {
		log.Fatalln("Error when loading location " + tz)
	}

	currentTime := time.Now().In(location)
	thenTime := currentTime.Add(1*time.Minute + 30*time.Second)
	now := currentTime.Format("2006-01-02 15:04:05")
	then := thenTime.Format("2006-01-02 15:04:05")

	return now, then
}

func get(url string, token string, checksum string, date string, expiration string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "fugw-edunext.fpt.edu.vn")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("origin", "https://fu-edunext.fpt.edu.vn")
	req.Header.Set("x-checksum", checksum)
	req.Header.Set("x-date", date)
	req.Header.Set("x-expiration", expiration)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyText)
}

func post(url string, token string, checksum string, date string, expiration string, body map[string]any) string {
	client := &http.Client{}
	bodyData, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyData))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("authority", "fugw-edunext.fpt.edu.vn")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("origin", "https://fu-edunext.fpt.edu.vn")
	req.Header.Set("X-Checksum", checksum)
	req.Header.Set("X-Date", date)
	req.Header.Set("X-Expiration", expiration)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyText)
}

func main() {
	const KEY = "331e4b21-a470-4833-89e0-db9ba605e8c7"
	const JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRodWFucG1jZTE1MDI1MUBmcHQuZWR1LnZuIiwibmFtZSI6IlBoYW4gTWluaCBUaHXhuq1uIiwicm9sZSI6NCwidXNlcklkIjozMTEwOCwiX2lkIjoiNjQ1N2FlMjU4ZGUyN2MwZDkxYjMyNjlhIiwiY2FtcHVzQ29kZSI6ImN0IiwiaWF0IjoxNjk1ODI0OTkwLCJleHAiOjE2OTYwODQxOTB9.t8K8BAMNNN9igE4V5nGU2OVAc92_aaAA6eX1pcGlMrs"
	const api = "https://fugw-edunext.fpt.edu.vn/api/comment/up-votes"

	// endpoint := getApiEndpoint(api)
	// fmt.Println(endpoint)

	now, then := getTime("Asia/Ho_Chi_Minh")

	payload := fmt.Sprintf("%s_%s_%s", "/api/comment/up-votes", now, then)

	checksum := hmacsha256(payload, KEY)

	body := map[string]any{
		"typeFilter": 1,
		"cqId":       80677,
		"commentId":  "64f7ee12687b97ebb4d605e3",
		"typeStar":   2,
		"star":       3,
	}

	data := post(api, JWT_TOKEN, checksum, now, then, body)

	fmt.Println(data)

}
