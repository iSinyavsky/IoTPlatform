package test

import (
	"fmt"
	"iot-project/tools"
	"net/http"
)

func init() {
	http.Handle("/api", tools.ModifyHTTPCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Yeah")
	})))
}
