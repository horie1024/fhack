package iqon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type IQON struct {
	Results []Results
}

type Results struct {
	Title      string
	Brand_name string
	Link       string
	Desc_long  string
	Price      string
	Image_link string
	Images     Images
}

type Images struct {
	L_image string
	M_image string
	S_image string
}

// IQON APIにリクエスト
func FetchIQON() IQON {

	resp, _ := http.Get("http://api.thefashionhack.com/item/?category_id1=10&page=1&limit=5&score_sort=1&instock_flag=1")
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var data IQON
	err := json.Unmarshal(byteArray, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
