package main

import (
	"net/http"
	"regexp"
	"strings"
)

var rt *regexp.Regexp

func init() {
	rt = regexp.MustCompile(`^\w+$`)
}

func sanitizeArgs(args []string) {
	//arguments: table_name/sort_key/sort_value
	if !rt.MatchString(args[0]) || len(args) > 1 && !rt.MatchString(args[1]) {
		panic(http.StatusBadRequest)
	}
	qry := `SELECT 1 FROM ` + args[0] + ` LIMIT 1`
	rows, err := db.Query(qry)
	if err != nil {
		rows.Close()
		panic(http.StatusNotFound)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case int:
				http.Error(w, http.StatusText(e.(int)), e.(int))
			default:
				http.Error(w, trace(e.(error).Error()).Error(), http.StatusInternalServerError)
			}
		}
	}()
	args := strings.Split(r.URL.Path[1:], "/")
	sanitizeArgs(args)
	//table := args[0]
	switch r.Method {
	case "GET":
		//qry := `SELECT * FROM '%s'`
	case "HEAD":
		//SELECT COUNT
	case "POST":
		if !cf.ALLOW_WRITE {
			panic(http.StatusForbidden)
		}
		//INSERT IGNORE INTO
	case "PATCH":
		if !cf.ALLOW_WRITE {
			panic(http.StatusForbidden)
		}
		//REPLACE INTO
	case "PUT":
		if !cf.ALLOW_WRITE || !cf.ALLOW_PUT {
			panic(http.StatusForbidden)
		}
		//FLUSH + INSERT
	default:
		panic(http.StatusMethodNotAllowed)
	}
}
