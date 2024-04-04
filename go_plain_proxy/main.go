package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	localPort := 3333
	remoteHost := "localhost"
	remotePort := 8080

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		remote, err := http.NewRequest(r.Method, fmt.Sprintf("http://%s:%d%s", remoteHost, remotePort, r.URL.Path), r.Body)
		if err != nil {
			fmt.Println("Error creating request:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		remote.Header = r.Header // Copy headers (excluding connection-related ones)

		client := &http.Client{}
		resp, err := client.Do(remote)
		if err != nil {
			fmt.Println("Error forwarding request:", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	})

	fmt.Println("Listening on port", localPort)
	http.ListenAndServe(fmt.Sprintf(":%d", localPort), nil)
}
