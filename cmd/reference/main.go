package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var Capabilities = []string{""}

type uniCorn struct {
	Name         string
	Capabilities []string
}

func main() {

	Capabilities = append(Capabilities, "super strong", "fullfill wishes", "fighting capabilities", "fly", "swim", "sing", "run", "cry", "change color", "talk", "dance", "code", "design", "drive", "walk", "talk chinese", "lazy")
	http.HandleFunc("/api/get-unicorn", GetUnicorn)
	http.ListenAndServe(":8888", nil)
}

func GetUnicorn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("processing new request..")
	fn, err := os.Open("petnames.txt")
	if err != nil {
		fmt.Println("please try later, unicorn factory unavailable")
		return
	}
	var names []string
	var scanner = bufio.NewScanner(fn)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	fa, err := os.Open("adj.txt")
	if err != nil {
		panic(err)
	}
	var adj []string
	scanner = bufio.NewScanner(fa)
	for scanner.Scan() {
		adj = append(adj, scanner.Text())
	}

	sleep_time := time.Duration(rand.Intn(1000)) * time.Millisecond
	values := r.URL.Query()
	amount, _ := strconv.Atoi(values.Get("amount"))

	items := []uniCorn{}
	for j := 1; j < amount; j++ {
		name := adj[rand.Intn(1345)] + "-" + names[rand.Intn(5800)]
		item := uniCorn{
			Name: name,
		}
		items = append(items, item)
		time.Sleep(sleep_time)

		for i := 0; i < 3; i++ {

			cap := Capabilities[rand.Intn(18)]
			item.Capabilities = append(item.Capabilities, cap)
		}
	}
	d, _ := json.Marshal(items)
	fmt.Println("Unicorn ready..")
	w.Write(d)
}
