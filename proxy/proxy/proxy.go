package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const proxyAddr string = "localhost:9000"

var (
	counter = 0
	hosts   = []string{
		"http://localhost:8080",
		"http://localhost:8081",
	}
)

func ProxyTwoReplicasRun() {
	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	numberProxy := counter % len(hosts)
	url := hosts[numberProxy] + r.RequestURI
	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	c := http.Client{}
	resp, err := c.Do(proxyReq)
	if err != nil {
		return
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(nil)
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(response)
	if err != nil {
		return
	}
	counter++
}
