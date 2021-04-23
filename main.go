package main

import (
	"log"
	"net/http"
)

func main()  {
	proxy := NewProxy()
	err := proxy.Connect("192.168.1.37:8083")
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if err := proxy.Close(); err != nil {
			log.Println(err)
		}
	}()

	http.HandleFunc("/", proxy.HTTPHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}