package httpreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func printRequest(r *http.Request) {
	fmt.Println("\n#", time.Now().Format(time.RFC3339))

	fmt.Printf("\n```\n")
	fmt.Printf("%v %v %v\n\n", r.Method, r.URL.String(), r.Proto)

	printHeader(r.Header)
	fmt.Println("")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Print("print request: failed to read the body: ", err)
		return
	}
	printBody(body)

	if r.Trailer != nil {
		fmt.Println("")
		printHeader(r.Trailer)
	}

	fmt.Printf("```\n")
}

func printHeader(header http.Header) {
	for h, vals := range header {
		fmt.Printf("%v: %v\n", h, strings.Join(vals, ", "))
	}
}

func printBody(body []byte) {
	pretty := bytes.Buffer{}
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		fmt.Println("print body: error: ", err)
		return
	}
	fmt.Println(pretty.String())
}

func printDoResponse(r *http.Response, err error) {
	if err != nil {
		fmt.Printf("request error: %v\n", err)
		return
	}
	defer r.Body.Close()

	fmt.Printf("\n%v %v\n\n", r.Proto, r.Status)
	printHeader(r.Header)
	fmt.Println("")
	if body, err := io.ReadAll(r.Body); err != nil {
		fmt.Print("print request: failed to read the body: ", err)
	} else {

		printBody(body)
	}
	fmt.Println("")
}
