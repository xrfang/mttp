package main

type Configuration struct {
	HTTP_PORT     string
	MODE          int //0=read-only; 1=write-only; 2=read-write
	PASSWORD      string
	READ_TIMEOUT  int
	WRITE_TIMEOUT int
}

var cf Configuration

func loadConfig(fn string) {
	//default values
	cf.HTTP_PORT = "6887"
	cf.READ_TIMEOUT = 60
	cf.WRITE_TIMEOUT = 60
	if fn != "" {
		assert(conf.ParseFile(fn, &cf))
	}
}
