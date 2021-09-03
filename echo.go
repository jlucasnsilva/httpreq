package httpreq

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// EchoServer ...
func EchoServer(addr string) error {
	fmt.Printf("\nrunning echo server @ http://localhost%v/\n", addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("unable to read request body: ", err)
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(b))

		printRequest(r)

		w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
		w.Write(b)
	})
	return http.ListenAndServe(addr, nil)
}
