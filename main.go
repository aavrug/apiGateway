package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	cfg "github.com/aavrug/apiGateway/config"
	// "strings"
)

var configfile = flag.String("config", "config.json", "config file path")

// isAuthorized checks authorization and returns response
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Cookie"] != nil {
			token := r.Header["Cookie"][0]

			if token != "" {
				endpoint(w, r)
				// More logic will be here
			} else {
				JSONError(w, "You are not authorized!", 401)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

// validateDataType checks authorization and returns response
func validateDataType(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Cookie"] != nil {
			token := r.Header["Cookie"][0]

			if token != "" {
				fmt.Println(token)
			} else {
				JSONError(w, "You are not authorized!", 401)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

// JSONError returns error copde
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	log.Printf("Status code: %v Error message: %v", code, err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func authRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Println("This is a get request")
	} else if r.Method == http.MethodPost {
		fmt.Println("This is a post request")
	} else if r.Method == http.MethodPut {
		fmt.Println("This is a " + r.Method + " request")
	} else if r.Method == http.MethodDelete {
		fmt.Println("Delete request")
	}
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", rBody)
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func todoRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		resp := sendRequest(w, r)
		fmt.Fprintf(w, resp)
		// fmt.Println("This is a get request")
	} else if r.Method == http.MethodPost {
		sendRequest(w, r)
		// fmt.Println("This is a post request")
	} else if r.Method == http.MethodPut {
		fmt.Println("This is a " + r.Method + " request")
	} else if r.Method == http.MethodDelete {
		fmt.Println("Delete request")
	}
}

func sendRequest(w http.ResponseWriter, r *http.Request) string {
	url := r.URL
	// config, err := getConfig()
	// fmt.Println(config)
	url.Host = "example.com"
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	proxyReq, err := http.NewRequest(r.Method, "https:"+url.String(), tee)
	if err != nil {
		fmt.Println(err)
	}

	// proxyReq.Header.Set("Host", r.Host)
	// proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	// fmt.Println(r.Header)
	// for header, values := range r.Header {
	// 	for _, value := range values {
	// 		proxyReq.Header.Add(header, value)
	// 	}
	// }

	client := &http.Client{}
	proxyRes, err := client.Do(proxyReq)
	fmt.Println(buf.String())

	if err != nil {
		fmt.Println(err)
	}

	defer proxyRes.Body.Close()
	// fmt.Println(proxyRes.Header)
	body, err := ioutil.ReadAll(proxyRes.Body)

	if err != nil {
		fmt.Println(err)
	}

	// var respBuf bytes.Buffer
	// responseTee := io.TeeReader(proxyRes.Body, &respBuf)
	// io.Copy(w, responseTee)
	// fmt.Println(respBuf.String())

	// fmt.Println("Response start")
	// fmt.Println(string(body))
	// fmt.Println("Response end")

	return string(body)
}

func validateRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	JSONError(w, "Invalid request!", 400)
}

func handleRequests() {
	http.HandleFunc("/", validateRequest)
	http.Handle("/auth", isAuthorized(authRequestHandler))
	http.Handle("/todo", isAuthorized(todoRequestHandler))
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func main() {
	if err := cfg.LoadConfig(*configfile); err != nil {
		log.Print("load condig file failed : ", err.Error())
		return
	}
	// config := cfg.GetConfig()
	// fmt.Printf("%#v\n", users)

	// fmt.Println(config["todo"].Host)
	log.Print("start init default server", cfg.GetConfig())
	// handleRequests()
	// log.Print("end init default server", cfg.GetConfig().EtcdConfig)
}
