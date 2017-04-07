package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func getJS() []byte {
	return []byte(`
		<script>
			var hash = location.hash;
			if (hash.startsWith("#")) {
				window.location = "http://localhost:49999/?"+hash.slice(1);
			}
		</script>
	`)
}

func getToken(port int) (string, error) {

	done := make(chan string)

	// This server waits for the redirect coming back from API server, populates
	// reqErr and returns the token from that request, and then stops itself.
	srv := &http.Server{
		Addr: fmt.Sprintf("localhost:%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This is to handle fragment parsing in implicit code flow
			if r.RequestURI == "/" {
				w.Write(getJS())
				return
			}

			if r.Method != "GET" {
				reqErr = errors.New("The server made a bad request: Only GET is allowed")
			}

			token := r.URL.Query().Get("access_token")
			if token != "" {
				done <- token
			}
		}),
	}
	go srv.ListenAndServe()

	token := <-done

	err := srv.Close()
	if err != nil {
		return token, errors.Wrap(err, "Error shutting down server")
	}
	return token, nil
}
