package server

import (
	"fmt"
	"net/http"
)

func ServeTls(handler http.Handler, cert, key string) error {
	fmt.Println("Server is running on port *:8443")
	return http.ListenAndServeTLS(":8443", cert, key, handler)
}

func Serve(tokenEndpointHandler http.Handler) error {
	fmt.Println("Server is running on port *:8080")
	return http.ListenAndServe(":8080", tokenEndpointHandler)
}
