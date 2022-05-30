package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func handleIncomingCall(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("headers: %v\n", r.Header)

	_, err := io.Copy(os.Stdout, r.Body)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	var portFlag = flag.Int("p", 8080, "port on which to serve http")
	flag.Parse()

	var port = fmt.Sprintf(":%v", strconv.Itoa(*portFlag))

	log.Println("server started on port", port)
	http.HandleFunc("/incoming_call", handleIncomingCall)
	log.Fatal(http.ListenAndServe(port, nil))
}
