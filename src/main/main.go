package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

type IQON struct {
	Info    info
	Results []results
}

type info struct {
	Total        int
	Return_count int
	Offset       int
	Page         int
	Total_page   int
}

type results struct {
	Title string
}

func Index(w http.ResponseWriter, r *http.Request) {

	resp, _ := http.Get("http://api.thefashionhack.com/item/?category_id1=10&page=1&limit=1&score_sort=1&instock_flag=1")
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var f IQON
	err := json.Unmarshal(byteArray, &f)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(f)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprint(w, string(data))
}

func main() {

	l, _ := net.Listen("tcp", ":9000")

	http.HandleFunc("/api", Index)

	fcgi.Serve(l, nil)
}
