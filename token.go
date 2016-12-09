package main

import (
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"net"
	"net/http"
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

func getToken(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	// This server waits for the redirect coming back from API server, populates
	// token and reqErr from that request, and then stops itself.
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This is to handle fragment parsing in implicit code flow
			//spew.Dump(r)
			if r.RequestURI == "/" {
				w.Write(getJS())
				return
			}

			// Stop listening once we've gotten a request.
			listener.Close()
			if r.Method != "GET" {
				reqErr = errors.New("The server made a bad request: Only GET is allowed")
			}

			token = r.URL.Query().Get("access_token")
			if token == "" {
				reqErr = errors.New("Missing 'token_type' parameter from server.")
			}

			wg.Done()
		}),
	}
	wg.Add(1)
	go srv.Serve(listener)
	return nil
}
