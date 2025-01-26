// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"
)

func main() {
	item_id := 1
	price := 1.1

	start := time.Now().UnixNano()

	var data = []byte(fmt.Sprintf(`{
		"ids": [
			%d
		],
		"partner": "359688187",
		"token": "YAwTGWnG",
		"max_price": %.2f,
		"custom_id": "%d"
	}`, item_id, price, item_id))

	// fmt.Println(string(data), item_id)

	// proxyUrl, err := url.Parse("http://194.99.27.66:8085")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	req, err := http.NewRequest(http.MethodPost, "https://api.lis-skins.com/v1/market/buy", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	// client := &http.Client{Timeout: 20 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, err := client.Do(req)

	fmt.Println("Buying finished for:", float64(time.Now().UnixNano()-start)/1000000000, "time now:", time.Now())

	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(fmt.Printf("status code error: %d %s", resp.StatusCode, resp.Status))
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	if resp.StatusCode == 201 {

		var answer map[string]interface{}

		err = json.Unmarshal([]byte(string(data)), &answer)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println(resp.StatusCode, answer["success"], reflect.TypeOf(answer["success"]))
	}
}
