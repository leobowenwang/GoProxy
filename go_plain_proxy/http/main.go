package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleProxy(w http.ResponseWriter, r *http.Request, remoteHost string, remotePort int) {
	targetURL := fmt.Sprintf("http://%s:%d", remoteHost, remotePort)

	req, err := http.NewRequest(r.Method, targetURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for header, values := range r.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleProxy(w, r, "localhost", 8080)
	})
	log.Fatal(http.ListenAndServe(":3333", nil))
}
