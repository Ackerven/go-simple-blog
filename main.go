package main

import (
	"fmt"
	"log"
	"net/http"
	"simple-blog/utils"
)

func printf(w http.ResponseWriter, r *http.Request)  {
	var str string
	str = fmt.Sprintf("%s: %s", "mysql", utils.DbName)
	fmt.Fprintf(w, str)
}

func main() {
	http.HandleFunc("/", printf)
	log.Fatal(http.ListenAndServe(utils.HttpPort, nil))
}