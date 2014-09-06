package main

import (
	"../lib/iqon"
	"../lib/place"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

type Response struct {
	Iqon  iqon.IQON
	Place []place.PlaceData
}

func Index(w http.ResponseWriter, r *http.Request) {

	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")

	//35.658643,139.7006439

	fmt.Println(lat)
	fmt.Println(lng)

	f := iqon.FetchIQON()

	//fmt.Println(f)

	loc := place.LOC{
		Lat: lat,
		Lng: lng,
	}

	// 現在位置と店舗の位置を計算する
	placeDataArray := place.Calc(place.FetchPlace(f, loc), loc)

	data := Response{
		Iqon:  f,
		Place: placeDataArray,
	}

	fmt.Println(data)

	response, err := json.Marshal(data)
	if err != nil {
		log.Fatal("josn encode error")
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprint(w, string(response))
}

func main() {

	l, _ := net.Listen("tcp", ":9000")

	http.HandleFunc("/api", Index)
	//http.HandleFunc("/") 送られてきた位置情報とredisにある位置情報を比較

	fcgi.Serve(l, nil)
}
