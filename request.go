package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Message_test struct {
	Success  bool                                `json:'success'`
	Currency string                              `json:'currency'`
	Data     map[string][]map[string]interface{} `json:'data'`
	// Data     []map[string]interface{} `json:'data'`
}

type Token_test struct {
	Data map[string]string `json:'data'`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var db_test *pgxpool.Pool

func test1() {
	for i := 0; i < 10; i++ {

		start := time.Now().UnixNano()

		resp, err := http.Get("https://market.csgo.com/api/v2/search-list-items-by-hash-name-all?key=RYmNCt39M717X78b8YIU6zm48Ss8QCf&list_hash_name[]=%E2%98%85%20Paracord%20Knife%20|%20Stained%20(Battle-Scarred)")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("status code error: %d %s", resp.StatusCode, resp.Status)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(string(data), reflect.TypeOf(data))

		var m Message_test

		err = json.Unmarshal([]byte(string(data)), &m)
		if err != nil {
			log.Fatal(err)
			return
		}

		for key, value := range m.Data {
			fmt.Println(key)
			// fmt.Println(value)

			for _, value := range value {
				fmt.Println(value["id"], reflect.TypeOf(value["id"]))

				id, ok := value["id"].(float64)
				if ok {
					fmt.Println(int(id) == 6084843520)
				}

			}
		}

		fmt.Println(float64(time.Now().UnixNano()-start) / 1000000000)

	}
}

func test2() {

	req, err := http.NewRequest(http.MethodGet, "https://api.lis-skins.com/v1/user/get-ws-token", nil)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer cfbb3eb6-dcfc-4649-9130-10e9bd7c8143")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	var answer Token_test

	err = json.Unmarshal([]byte(string(data)), &answer)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(answer.Data["token"])

}

func test3() {

	for range 5 {

		start := time.Now().UnixNano()

		var price float64
		var date1 int
		err := db_test.QueryRow(context.Background(), "SELECT price, NOW()::date - last_update_price::date FROM tm.items WHERE name=$1", "StatTrak™ AK-47 | Asiimov (Field-Tested)").Scan(&price, &date1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(float64(time.Now().UnixNano()-start) / 1000000000)

		fmt.Println(price, date1)

	}

}

type Buy_answer struct {
	Success  bool                                `json:'success'`
	Currency string                              `json:'currency'`
	Data     map[string][]map[string]interface{} `json:'data'`
	// Data     []map[string]interface{} `json:'data'`
}

func Buy(account_id int, user string, item_id int, price float64) {

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

	proxyUrl, err := url.Parse("http://194.99.27.66:8085")
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.lis-skins.com/v1/market/buy", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Add("Authorization", "Bearer 3f48afa0-0102-4f32-b235-dbdee027d688")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
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

func Buy_test(account_id int, user string, item_id int, price float64) {

	for range 2 {
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

		req, err := http.NewRequest(http.MethodPost, "https://api.lis-skins.com/v1/market/buy", bytes.NewBuffer(data))
		if err != nil {
			log.Fatal(err)
			return
		}

		req.Header.Add("Authorization", "Bearer 3f48afa0-0102-4f32-b235-dbdee027d688")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		fmt.Println("Buying finish for:", float64(time.Now().UnixNano()-start)/1000000000, "time now:", time.Now())

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

		time.Sleep(5 * time.Second)
	}
}

func mai1n() {

	// conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/trade")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(context.Background())

	// poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	// poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/trade?pool_max_conns=40")
	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/trade")
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	// poolConfig.MinConns = 0
	// poolConfig.MaxConns = 60

	db_test, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	// Buy(1, "Artem", 122699506, 100000)

	for range 2 {
		Buy_test(1, "1", 1, 1)
	}
}
