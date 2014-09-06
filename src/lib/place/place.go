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

type PlaceArray struct {
	Results results
	Flag    bool
}

type LOC struct {
	Lat string
	Lng string
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

func FetchPlace(iqon iqon.IQON, loc LOC) []results {

	// place api リクエスト
	client := &http.Client{}
	values := url.Values{}
	values.Add("key", "AIzaSyAHWuu8QLFiD9P6zI1q8CHD_-5RhckWUs4")
	values.Add("sensor", "false") //あとでtrue
	values.Add("radius", "1000")
	values.Add("keyword", iqon.Results[0].Brand_name)

	//fmt.Println(values.Encode())

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

	//fmt.Println(data)

	var d []results

	for i, j := range data.Results {

		for _, n := range j.Types {

			// 何でフィルタリングするかは調整必要かも
			if n == "department_store" {

				//fmt.Println(i, n)
				//fmt.Println(j)
				d = append(d, data.Results[i])
			}
		}
	}

	return d
}

func Calc(p []results, loc LOC) []PlaceArray {

	var distanceArray []PlaceArray

	for i, j := range p {

		fmt.Println(j.Geometry.Location.Lat)
		fmt.Println(j.Geometry.Location.Lng)

		floatLat, _ := strconv.ParseFloat(loc.Lat, 32)
		floatLng, _ := strconv.ParseFloat(loc.Lng, 32)

		latDiff := math.Pow((j.Geometry.Location.Lat-
			floatLat)/0.0111, 2)

		lngDiff := math.Pow((j.Geometry.Location.Lng-floatLng)/0.0091, 2)

		a := math.Sqrt(latDiff + lngDiff)
		fmt.Println(a)

		distance := a * 1000
		if distance <= 50 {

			b := PlaceArray{
				Results: p[i],
				Flag:    true,
			}

			distanceArray = append(distanceArray, b)
		} else {

			b := PlaceArray{
				Results: p[i],
				Flag:    false,
			}

			distanceArray = append(distanceArray, b)
		}
	}

	return distanceArray
}
