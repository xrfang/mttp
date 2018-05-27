package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	ver := flag.Bool("version", false, "show version info")
	conf := flag.String("conf", "", "configuration file")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *conf == "" {
		fmt.Println("missing configuration (-conf)")
		return
	}
	loadConfig(*conf)
	res := Encrypt([]byte("hello world!"))
	fmt.Println(base64.StdEncoding.EncodeToString(res))
	rev := Decrypt(res)
	fmt.Println(string(rev))
	return
	http.HandleFunc("/", handler)
	svr := http.Server{
		Addr:         ":" + cf.HTTP_PORT,
		ReadTimeout:  time.Duration(cf.READ_TIMEOUT) * time.Second,
		WriteTimeout: time.Duration(cf.WRITE_TIMEOUT) * time.Second,
	}
	svr.ListenAndServe()
}
