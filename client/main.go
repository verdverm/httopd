package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	numClient   = flag.Int("c", 4, "number of client routines")
	meanReqTime = flag.Float64("mean", 3.0, "mean time between request per client")
	sdevReqTime = flag.Float64("sdev", 1.0, "mean time between request per client")
	hostStr     = os.Getenv("HOSTURL")
	mean, sdev  float64
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	mean, sdev = *meanReqTime, *sdevReqTime

	fmt.Println("http client simulator")

	for i := 1; i < *numClient; i++ {
		go run(i)
	}
	run(0)
}

func run(id int) {
	for {
		url := genUrl()
		fmt.Println(id, url)
		getUrl(url)
		sleep()
	}
}

func genUrl() string {
	// generate random number in the range [1:5]
	// where 5 has half the probability of the other values
	// (5 is a 404 page)
	r := rand.Int()
	i := ((r % 9) + 2) / 2
	url := "http://localhost:8080/page" + fmt.Sprintf("%d", i)
	// add extra level on page
	if r%2 == 0 {
		url += "/asubpage"
	}
	if r%3 == 0 {
		url += "/asubsubpage"
	}
	return url
}

func getUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Get Error: ", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll Error: ", err)
	}
}

func sleep() {
	f := rand.NormFloat64()*sdev + mean
	i := int(math.Floor(f + 0.5))
	dur := time.Second * time.Duration(i)
	time.Sleep(dur)
}
