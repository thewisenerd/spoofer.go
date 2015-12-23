/*
 * spoofer.go
 * Copyright 2015 thewisenerd <thewisenerd@protonmail.com>
 *
 * Use of this source code is governed by a GNU GPL v2.0
 * license that can be found in the LICENSE file.
 *
 */

package main

import (
	"fmt"
	"strings"
	"strconv"
	"io"
	"io/ioutil"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	//TODO: write documentation.
	io.WriteString(w, http.StatusText(http.StatusOK))
}

func spoof(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	// url param not set
	if len(q) == 0 {
		// malformed request
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "missing url param")
		return
	}

	url := r.URL.Query().Get("url")

	// url param empty
	if url == "" {
		// malformed request
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "missing url param")
		return
	}

	// referer param set
	if r.URL.Query().Get("referer") != "" {
		r.Header.Set("referer", r.URL.Query().Get("referer"))
	}

	// invalidate cache from client
	// set no-cache
	r.Header.Set("Pragma", "no-cache"); //HTTP 1.0
	r.Header.Set("Cache-Control", "no-cache, must-revalidate"); //HTTP1.1

	// remove If-Modified-Since
	_, ok := r.Header["If-Modified-Since"];
	if ok {
		delete(r.Header, "If-Modified-Since");
	}

	fmt.Printf("%s: ", url)

	// initialize client
	// TODO: add timeout
	client := &http.Client{};

	// TODO: support non-GET HTTP methods
	// TODO: forward r.Method
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// something went wrong
		fmt.Printf("%d\n", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	// set req header
	req.Header = r.Header

	// get response
	resp, err := client.Do(req)
	if err != nil {
		// something went wrong
		fmt.Printf("%d\n", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	fmt.Printf("%d\n", status)
	if (status != http.StatusOK) {
		w.WriteHeader(status)
		io.WriteString(w, http.StatusText(status))
		return
	}

	// read resp.body
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// something went wrong
		fmt.Printf("%d\n", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	// fw headers
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ", "))
	}

	// invalidate cache for client
	w.Header().Set("Pragma", "no-cache"); //HTTP 1.0
	w.Header().Set("Cache-Control", "no-cache, must-revalidate"); //HTTP1.1

	// fix Content-Length if necessary
	if (len(contents) != 0) && (w.Header().Get("Content-Length") == "0") {
		w.Header().Set("Content-Length", strconv.Itoa(len(contents)))
	}

	// fw content
	w.Write(contents)
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = hello
	mux["/spoof"] = spoof

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.Path]; ok {
		h(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, http.StatusText(http.StatusNotFound))
}
