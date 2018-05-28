package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ver := flag.Bool("version", false, "show version info")
	conf := flag.String("conf", "", "configuration file")
	dec := flag.String("decrypt", "", "decrypt given file")
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
	if *dec != "" {
		var buf bytes.Buffer
		if *dec == "-" {
			*dec = "/dev/stdin"
		}
		f, err := os.Open(*dec)
		assert(err)
		defer f.Close()
		_, err = io.Copy(&buf, f)
		assert(err)
		data := bytes.NewBuffer(Decrypt(buf.Bytes()))
		zr, err := gzip.NewReader(data)
		assert(err)
		_, err = io.Copy(os.Stdout, zr)
		assert(err)
		return
	}
	http.HandleFunc("/", handler)
	svr := http.Server{
		Addr:         ":" + cf.HTTP_PORT,
		ReadTimeout:  time.Duration(cf.READ_TIMEOUT) * time.Second,
		WriteTimeout: time.Duration(cf.WRITE_TIMEOUT) * time.Second,
	}
	assert(svr.ListenAndServe())
}
