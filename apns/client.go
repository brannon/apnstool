// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

const (
	DeviceEndpointFormat = "https://%s/3/device/%s"

	ProductionEndpoint = "api.push.apple.com"
	SandboxEndpoint    = "api.sandbox.push.apple.com"
)

type Headers map[string]string

type Client interface {
	ConfigureCertificateAuth(cert tls.Certificate)
	ConfigureEndpoint(endpoint string)
	ConfigureTokenAuth(token string)
	EnableLogging(writer io.Writer)
	Send(deviceToken string, headers Headers, content []byte) (*SendResult, error)
}

type client struct {
	bearerToken string
	certificate tls.Certificate
	endpoint    string
	logWriter   io.Writer
}

func NewClient() Client {
	return &client{
		bearerToken: "",
		certificate: tls.Certificate{},
		endpoint:    ProductionEndpoint,
		logWriter:   nil,
	}
}

func (c *client) ConfigureCertificateAuth(cert tls.Certificate) {
	c.certificate = cert
}

func (c *client) ConfigureEndpoint(endpoint string) {
	c.endpoint = endpoint
}

func (c *client) ConfigureTokenAuth(token string) {
	c.bearerToken = token
}

func (c *client) EnableLogging(writer io.Writer) {
	c.logWriter = writer
}

func (c *client) Send(deviceToken string, headers Headers, content []byte) (*SendResult, error) {
	deviceUrl, err := url.Parse(fmt.Sprintf(DeviceEndpointFormat, c.endpoint, deviceToken))
	if err != nil {
		return nil, err
	}

	var transport http.RoundTripper

	if c.certificate.PrivateKey != nil {
		c.log("* Using client certificate\n")
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{
				c.certificate,
			},
		}

		tlsTransport := &http.Transport{
			TLSClientConfig:   tlsConfig,
			ForceAttemptHTTP2: true,
		}

		http2.ConfigureTransport(tlsTransport)
		transport = tlsTransport
	} else {
		transport = http.DefaultTransport
	}

	client := http.Client{Transport: transport}

	req := &http.Request{
		Method:     "POST",
		URL:        deviceUrl,
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header:     make(http.Header),
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
	}

	c.log("* Sending request:\n")
	c.logf("> %s %s\n", req.Method, req.URL.String())
	for name, _ := range req.Header {
		c.logf("> %s: %s\n", name, req.Header.Get(name))
	}
	c.logf("> %s\n", content)

	req.Body = ioutil.NopCloser(bytes.NewReader(content))

	res, err := client.Do(req)
	if err != nil {
		c.logf("* Error sending request: %s\n", err)
		return nil, err
	}

	c.logf("* Received response:\n")
	c.logf("< %s\n", res.Status)
	for name, _ := range res.Header {
		c.logf("< %s: %s\n", name, res.Header.Get(name))
	}

	responseContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.logf("* Error reading response body: %s\n", err)
		return nil, err
	}

	c.logf("< %s\n", responseContent)

	result := &SendResult{
		content:    responseContent,
		headers:    res.Header,
		StatusCode: res.StatusCode,
	}

	c.log("* Done\n")

	return result, nil
}

func (c *client) log(text string) {
	if c.logWriter != nil {
		_, _ = io.WriteString(c.logWriter, text)
	}
}

func (c *client) logf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
}

type SendResult struct {
	content    []byte
	headers    http.Header
	StatusCode int
}

func (r *SendResult) ErrorReason() string {
	if r.content != nil && len(r.content) > 0 {
		var object map[string]interface{}

		_ = json.Unmarshal(r.content, &object)

		if reason, ok := object["reason"].(string); ok {
			return reason
		}
	}
	return ""
}

func (r *SendResult) Id() string {
	return r.headers.Get("apns-id")
}

func (r *SendResult) Success() bool {
	return r.StatusCode == 200
}
