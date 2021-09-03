package httpreq

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	// Group ...
	Group struct {
		Config        Config            `yaml:"given"`
		EndpointCalls []EndpointCall    `yaml:"then"`
		Macros        map[string]string `yaml:"where,omitempty"`
	}

	// Config ...
	Config struct {
		HTTPS  bool              `yaml:"https"`
		Host   string            `yaml:"host"`
		Port   string            `yaml:"port"`
		Header map[string]string `yaml:"headers"`
	}

	// EndpointCall ...
	EndpointCall struct {
		Call    string            `yaml:"call"`
		ID      string            `yaml:"id"`
		Header  map[string]string `yaml:"with_headers,omitempty"`
		Payload string            `yaml:"and_send"`
	}

	requestOption struct {
		Request *http.Request
		Error   error
	}
)

// Execute ...
func (g *Group) Execute(ids ...string) {
	if len(ids) < 1 {
		exexuteAll(g)
		return
	}

	for _, id := range ids {
		executeOne(g, id)
	}
}

func exexuteAll(g *Group) {
	requests, err := createRequests(g)
	failOnError("execute: all:", err)

	for reqID, opt := range requests {
		if opt.Error != nil {
			fmt.Println("execute: all:", err)
			continue
		}

		client := http.Client{Timeout: 5 * time.Second}
		fmt.Printf("executing '%v'\n\n", reqID)
		printDoResponse(client.Do(opt.Request))
	}
}

func executeOne(g *Group, requestID string) {
	requests, err := createRequests(g)
	failOnError("execute: one:", err)

	opt, ok := requests[requestID]
	if !ok {
		fmt.Printf("execute: one: request '%v' not found", requestID)
		return
	}

	if opt.Error != nil {
		fmt.Println("execute: one:", opt.Error)
		return
	}

	client := http.Client{Timeout: 5 * time.Second}
	printDoResponse(client.Do(opt.Request))
}

func createRequests(g *Group) (map[string]requestOption, error) {
	url := makeBaseURL(g.Config)
	m := make(map[string]requestOption)

	for _, c := range g.EndpointCalls {
		id := c.ID
		method, uri, err := splitCall(c.Call)
		if err != nil {
			e := fmt.Errorf("execute: create requests: '%v' creation error: %v", id, err)
			m[id] = requestOption{Error: e}
			continue
		}

		url := joinURL(url, uri)
		rdr := strings.NewReader(c.Payload)
		req, err := http.NewRequest(method, url, rdr)
		if err != nil {
			e := fmt.Errorf("execute: create requests: '%v' creation error: %v", id, err)
			m[id] = requestOption{Error: e}
			continue
		}

		addHeader(req.Header, g.Config.Header, c.Header)
		setContentType(req.Header)

		m[id] = requestOption{Request: req}
	}
	return m, nil
}

func splitCall(call string) (method, uri string, err error) {
	parts := strings.Split(call, " ")
	if len(parts) != 2 {
		return "", "", errors.New("a call should have two elements: method and URI")
	}

	method = parts[0]
	if _, ok := httpMethods[method]; !ok {
		return "", "", fmt.Errorf("'%v' is not a supported HTTP method", method)
	}

	uri = parts[1]
	if len(uri) < 1 {
		return "", "", errors.New("the URI should have at least one character ('/' for instance)")
	}
	return method, uri, nil
}

func makeBaseURL(cfg Config) string {
	protocol := "http"
	if cfg.HTTPS {
		protocol = "https"
	}

	h := strings.TrimSpace(cfg.Host)
	p := strings.TrimSpace(cfg.Port)

	if p == "" {
		return fmt.Sprintf("%v://%v", protocol, h)
	}
	return fmt.Sprintf("%v://%v", protocol, net.JoinHostPort(h, p))
}

func joinURL(base, uri string) string {
	u := filepath.Clean(uri)
	if strings.HasPrefix(u, "/") {
		return base + u
	}
	return base + "/" + u
}

func setContentType(header http.Header) {
	for k := range header {
		if strings.ToLower(k) == "content-type" {
			return
		}
	}
	header.Set("Content-type", "application/json; charset=utf-8")
}

func addHeader(header http.Header, globalHeader, callHeader map[string]string) {
	for k, v := range globalHeader {
		for _, w := range strings.Split(v, ",") {
			header.Add(k, w)
		}
	}
	for k, v := range callHeader {
		for _, w := range strings.Split(v, ",") {
			header.Add(k, w)
		}
	}
}

func failOnError(prefix string, err error) {
	if err == nil {
		return
	}
	fmt.Println(prefix, err)
	os.Exit(1)
}
