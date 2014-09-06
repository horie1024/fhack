package place

import (
	"../iqon"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type LOC struct {
	Lat string
	Lng string
}

type PlaceData struct {
	Geometry Geometry
	Name     string
	Types    []string
	Place_id string
	Flag     bool
}

type Place struct {
	Results []results
}

type results struct {
	Geometry Geometry
	Name     string
	Types    []string
	Place_id string
}

type Geometry struct {
	Location Location
}

type Location struct {
	Lat float64
	Lng float64
}

func FetchPlace(iqon iqon.Results, loc LOC) []results {

	// place api リクエスト
	client := &http.Client{}
	values := url.Values{}
	values.Add("key", "AIzaSyAHWuu8QLFiD9P6zI1q8CHD_-5RhckWUs4")
	values.Add("sensor", "false") //あとでtrue
	values.Add("radius", "3000")
	values.Add("keyword", iqon.Brand_name)

	// Request を生成
	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/place/nearbysearch/json", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.URL.RawQuery = values.Encode() + "&location=" + loc.Lat + "," + loc.Lng

	// リクエスト
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	// byteに変換
	placeByteArray, _ := ioutil.ReadAll(resp.Body)

	var data Place
	err = json.Unmarshal(placeByteArray, &data)
	if err != nil {
		log.Fatal(err)
	}

	var d []results
	for i, j := range data.Results {

		for _, n := range j.Types {

			// 何でフィルタリングするかは調整必要かも
			if n == "department_store" || n == "clothing_store" {
				//if n == "department_store" {

				d = append(d, data.Results[i])
			}
		}
	}

	return d
}

func Calc(p []results, loc LOC) []PlaceData {

	var distanceArray []PlaceData

	for i, j := range p {

		floatLat, _ := strconv.ParseFloat(loc.Lat, 32)
		floatLng, _ := strconv.ParseFloat(loc.Lng, 32)

		latDiff := math.Pow((j.Geometry.Location.Lat-
			floatLat)/0.0111, 2)

		lngDiff := math.Pow((j.Geometry.Location.Lng-floatLng)/0.0091, 2)

		a := math.Sqrt(latDiff + lngDiff)

		distance := a * 1000
		if distance <= 300 && len(distanceArray) <= 5 {

			b := PlaceData{
				Name:     p[i].Name,
				Types:    p[i].Types,
				Place_id: p[i].Place_id,
				Geometry: Geometry{
					Location{
						Lat: p[i].Geometry.Location.Lat,
						Lng: p[i].Geometry.Location.Lng,
					},
				},
				Flag: true,
			}

			distanceArray = append(distanceArray, b)
		}
		/* else {

			b := PlaceData{
				Name:     p[i].Name,
				Types:    p[i].Types,
				Place_id: p[i].Place_id,
				Geometry: Geometry{
					Location{
						Lat: p[i].Geometry.Location.Lat,
						Lng: p[i].Geometry.Location.Lng,
					},
				},
				Flag: false,
			}

			distanceArray = append(distanceArray, b)
		}*/
	}

	return distanceArray
}
