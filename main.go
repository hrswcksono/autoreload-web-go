package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status `json:"status"`
}

type StringStatus struct {
	WaterStatus string
	WindStatus  string
}

func statusToString(input Status) StringStatus {
	var output = StringStatus{}

	switch {
	case input.Water < 5:
		output.WaterStatus = "aman"
	case input.Water < 9 && input.Water > 5:
		output.WaterStatus = "siaga"
	case input.Water > 8:
		output.WaterStatus = "bahaya"
	}

	switch {
	case input.Wind < 6:
		output.WindStatus = "aman"
	case input.Wind > 6 && input.Wind < 16:
		output.WindStatus = "siaga"
	case input.Wind > 15:
		output.WindStatus = "bahaya"
	}
	return output
}

func updateData() {
	for {

		var data = Data{Status: Status{}}

		data.Status.Water = rand.Intn(99) + 1

		data.Status.Wind = rand.Intn(99) + 1

		b, err := json.MarshalIndent(&data, "", " ")

		if err != nil {
			log.Fatalln("error while marshalling json data  =>", err.Error())
		}

		err = ioutil.WriteFile("data.json", b, 0644)

		if err != nil {
			log.Fatalln("error while writing value to data.json file  =>", err.Error())
		}
		fmt.Println("menggungu 5 detik")
		time.Sleep(time.Second * 5)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		tpl, _ := template.ParseFiles("index.html")

		var data = Data{Status: Status{}}

		b, err := ioutil.ReadFile("data.json")

		if err != nil {
			fmt.Fprint(w, "error open json")
			return
		}

		err = json.Unmarshal(b, &data)

		if err != nil {
			fmt.Fprint(w, "error unmarshal")
			return
		}

		err = tpl.ExecuteTemplate(w, "index.html", statusToString(data.Status))

		if err != nil {
			fmt.Fprint(w, "error execute template")
			return
		}

	})

	http.ListenAndServe(":8080", nil)
}
