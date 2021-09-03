package httpreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/ttacon/chalk"
)

func printRequest(r *http.Request) {
	printTitle(time.Now().Format(time.RFC3339))

	printBlockDiv()
	printRequestSignature(r.Method, r.URL.String(), r.Proto)

	printHeader(r.Header)

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

	printBlockDiv()
}

func printBlockDiv() {
	if !isTerminal() {
		fmt.Println("```")
	}
}

func printRequestSignature(method, url, proto string) {
	if isTerminal() {
		printRequestSignatureTerminal(method, url, proto)
	} else {
		printRequestSignatureMarkdown(method, url, proto)
	}
	fmt.Println("")
}

func printRequestSignatureTerminal(method, url, proto string) {
	m := fmt.Sprintf("%v%v%v", chalk.Red, chalk.Bold.TextStyle(method), chalk.Reset)
	u := fmt.Sprintf("%v%v%v", chalk.Cyan, url, chalk.Reset)
	fmt.Println(m, u, proto)
}

func printRequestSignatureMarkdown(method, url, proto string) {
	fmt.Println(method, url, proto)
}

func printTitle(title string) {
	fmt.Println("")
	if isTerminal() {
		printTitleTerminal(title)
	} else {
		printTitleMarkdown(title)
	}
	fmt.Println("")
}

func printTitleMarkdown(title string) {
	fmt.Println("#", title)
}

func printTitleTerminal(title string) {
	upper := strings.ToUpper(title)
	bold := chalk.Bold.TextStyle(upper)
	fmt.Println(chalk.Underline.TextStyle(bold))
}

func printHeader(header http.Header) {
	if isTerminal() {
		printHeaderTerminal(header)
	} else {
		printHeaderMarkdown(header)
	}
	fmt.Println("")
}

func printHeaderTerminal(header http.Header) {
	for h, vals := range header {
		p := fmt.Sprintf("%v%v%v:", chalk.Magenta, h, chalk.Reset)
		fmt.Println(p, strings.Join(vals, ", "))
	}
}

func printHeaderMarkdown(header http.Header) {
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

	if b := pretty.String(); isTerminal() {
		printBodyTerminal(b)
	} else {
		printBodyMarkdown(b)
	}
}

func printBodyTerminal(body string) {
	fmt.Printf("%v%v%v\n", chalk.Cyan, body, chalk.Reset)
}

func printBodyMarkdown(body string) {
	fmt.Println(body)
}

func printDoResponse(r *http.Response, err error) {
	if err != nil {
		fmt.Printf("request error: %v\n", err)
		return
	}
	defer r.Body.Close()

	printResponseSignature(r.Proto, r.Status)
	printHeader(r.Header)
	if body, err := io.ReadAll(r.Body); err != nil {
		fmt.Println("print request: failed to read the body: ", err)
	} else {
		printBody(body)
	}
}

func printResponseSignature(proto, status string) {
	if isTerminal() {
		printResponseSignatureTerminal(proto, status)
	} else {
		printResponseSignatureMarkdown(proto, status)
	}
	fmt.Println("")
}

func printResponseSignatureTerminal(proto, status string) {
	s := fmt.Sprintf("%v%v%v", chalk.Red, chalk.Bold.TextStyle(status), chalk.Reset)
	fmt.Println(proto, s)
}

func printResponseSignatureMarkdown(proto, status string) {
	fmt.Println(proto, status)
}

func printID(id string) {
	fmt.Println("")
	if isTerminal() {
		printIDTerminal(id)
	} else {
		printIDMarkdown(id)
	}
	fmt.Println("")
}

func printIDTerminal(id string) {
	e := fmt.Sprintf("%v%v%v", chalk.Red, chalk.Bold.TextStyle("Executing:"), chalk.Reset)
	idx := chalk.Bold.TextStyle(id)
	fmt.Println(e, idx)
}

func printIDMarkdown(id string) {
	fmt.Println("Executing:", id)
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}
