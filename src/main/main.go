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

type IQON struct {
	Results []results
}

type results struct {
	Title      string
	Brand_name string
	Link       string
	Desc_long  string
	Price      string
	Image_link string
	Images     iqon.Images
	Place      []place.PlaceData
}

type response struct {
	Iqon IQON
}

func Index(w http.ResponseWriter, r *http.Request) {

	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")

	// もしクエリが無ければ渋谷駅をデフォルトにする
	if lat == "" {
		lat = "35.658643"
	}
	if lng == "" {
		lng = "139.7006439"
	}

	loc := place.LOC{
		Lat: lat,
		Lng: lng,
	}

	items := iqon.FetchIQON()

	var data []results
	for _, item := range items.Results {

		placeDataArray := place.Calc(place.FetchPlace(item, loc), loc)

		// もしplaceDataArrayがnilだったら空データで初期化
		if placeDataArray == nil {
			placeDataArray = append(placeDataArray, place.PlaceData{
				Name:     "",
				Types:    []string{},
				Place_id: "",
				Geometry: place.Geometry{
					place.Location{
						Lat: 0.0,
						Lng: 0.0,
					},
				},
				Flag: false,
			})
		}

		itemData := results{
			Title:      item.Title,
			Brand_name: item.Brand_name,
			Link:       item.Link,
			Desc_long:  item.Desc_long,
			Price:      item.Price,
			Image_link: item.Image_link,
			Images:     item.Images,
			Place:      placeDataArray,
		}

		data = append(data, itemData)
	}

	Iqon := response{
		Iqon: IQON{
			Results: data,
		},
	}

	response, err := json.Marshal(Iqon)
	if err != nil {
		log.Fatal("josn encode error")
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	fmt.Fprint(w, string(response))
}

func main() {

	l, _ := net.Listen("tcp", ":9000")

	http.HandleFunc("/api", Index)
	//http.HandleFunc("/temp") // 温度
	// ハーストで無いなら、温度でplaceのサーチ範囲変える

	fcgi.Serve(l, nil)
}
