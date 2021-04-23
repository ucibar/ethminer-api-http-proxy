package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	minerAddress := flag.String("miner", "", "Miner API IP:Port")
	serveAddress := flag.String("serve", ":8081", "HTTP Server API:Port")
	flag.Parse()

	if *minerAddress == "" {
		flag.PrintDefaults()
		return
	}

	proxy := NewProxy()
	err := proxy.Connect(*minerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if err := proxy.Close(); err != nil {
			log.Println(err)
		}
	}()

	http.HandleFunc("/", proxy.HTTPHandler)
	log.Fatal(http.ListenAndServe(*serveAddress, nil))
}
