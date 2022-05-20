package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	Origin     string
	token      string
	HTTPClient *http.Client
}

func NewClient(origin, basePath string) *Client {
	return &Client{
		BaseURL: origin + basePath,
		Origin:  origin,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) Post(path string, v interface{}, body interface{}) error {
	json, err := json.Marshal(body)
	if err != nil {
		LogFatalError("Failed to make object into JSON", err)
	}
	req := c.newRequest("POST", path, bytes.NewBuffer(json))
	return c.sendRequest(req, v)
}

func (c *Client) Get(path string, v interface{}) error {
	req := c.newRequest("GET", path, nil)
	return c.sendRequest(req, v)
}

func (c *Client) newRequest(method string, path string, body io.Reader) *http.Request {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		LogFatalError(fmt.Sprintf("Failed to create %s request for %q", method, url), err)
	}
	return req
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Add("Origin", "https://pservice-permit.giantleap.no")
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "no,nb")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/json, text/plain, */*; charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	if c.token != "" {
		req.Header.Set("X-Token", c.token)
	} else {
		req.Header.Set("X-Token", "null")
	}

	if *vvFalg {
		log.Println(req.Method, req.URL.String(), req.Header)
		if req.Body != nil {
			log.Println("Body:", req.Body)
		}
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if *vvFalg {
		log.Println(res.StatusCode, req.Method, req.URL.String())
	}

	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	if *vvFalg {
		j, _ := json.MarshalIndent(v, "", "  ")
		log.Println(res.StatusCode, req.Method, req.URL.String(), string(j))
	}

	return nil
}
