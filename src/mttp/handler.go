package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
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
	_, err := db.Query(qry)
	if err != nil {
		panic(http.StatusNotFound)
	}
}

func getPayload(r *http.Request) (data []map[string]interface{}) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r.Body)
	assert(err)
	zbuf := bytes.NewBuffer(Decrypt(buf.Bytes()))
	zr, err := gzip.NewReader(zbuf)
	assert(err)
	defer zr.Close()
	assert(json.NewDecoder(zr).Decode(&data))
	return
}

func prepSql(cmd, tbl string, payload []map[string]interface{}) (sql string, args []interface{}) {
	if len(payload) == 0 {
		panic(http.StatusNoContent)
	}
	var keys []string
	for k := range payload[0] {
		keys = append(keys, k)
	}
	ph := "(?" + strings.Repeat(`,?`, len(payload[0])-1) + ")"
	sql = fmt.Sprintf(`%s INTO %s (%s) VALUES `, cmd, tbl, strings.Join(keys, ","))
	sql += ph + strings.Repeat(","+ph, len(payload)-1)
	for _, p := range payload {
		for _, k := range keys {
			args = append(args, p[k])
		}
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case int:
				http.Error(w, http.StatusText(e.(int)), e.(int))
			default:
				http.Error(w, trace("%v", e).Error(), http.StatusInternalServerError)
			}
		}
	}()
	args := strings.Split(r.URL.Path[1:], "/")
	if len(args) < 1 || len(args) > 3 {
		panic(http.StatusBadRequest)
	}
	sanitizeArgs(args)
	switch r.Method {
	case "GET":
		var where, orderby, limit string
		var c []interface{}
		selFrom := `SELECT * FROM ` + args[0]
		cnt := r.URL.Query().Get("limit")
		if cnt == "" {
			limit = " LIMIT " + strconv.Itoa(cf.READ_LIMIT)
		} else {
			limit = " LIMIT " + cnt
		}
		if len(args) > 1 {
			orderby = " ORDER BY " + args[1]
		}
		if len(args) > 2 {
			where = " WHERE " + args[1] + ">?"
			c = []interface{}{args[2]}
		}
		qry := selFrom + where + orderby + limit
		rows, err := db.Query(qry, c...)
		assert(err)
		raw := FetchRows(rows)
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		assert(json.NewEncoder(zw).Encode(raw))
		zw.Close()
		output := Encrypt(buf.Bytes())
		_, err = w.Write(output)
		assert(err)
	case "HEAD":
		sel := []string{`COUNT(1) AS 'Content-Length'`}
		if len(args) > 1 {
			sel = append(sel, `MAX(`+args[1]+`) AS 'X-Latest'`)
		}
		qry := `SELECT ` + strings.Join(sel, ",") + ` FROM ` + args[0]
		rows, err := db.Query(qry)
		assert(err)
		for k, v := range FetchRows(rows)[0] {
			w.Header().Add(k, fmt.Sprintf("%v", v))
		}
	case "POST":
		if !cf.ALLOW_WRITE {
			panic(http.StatusForbidden)
		}
		payload := getPayload(r)
		sql, data := prepSql("INSERT IGNORE", args[0], payload)
		res, err := db.Exec(sql, data...)
		assert(err)
		ra, _ := res.RowsAffected()
		out := []string{strconv.Itoa(int(ra))}
		if len(args) > 1 {
			last := payload[len(payload)-1]
			out = append(out, fmt.Sprintf("%v", last[args[1]]))
		}
		fmt.Fprintln(w, strings.Join(out, ","))
	case "PATCH":
		if !cf.ALLOW_WRITE {
			panic(http.StatusForbidden)
		}
		payload := getPayload(r)
		sql, data := prepSql("REPLACE", args[0], payload)
		res, err := db.Exec(sql, data...)
		assert(err)
		ra, _ := res.RowsAffected()
		out := []string{strconv.Itoa(int(ra))}
		if len(args) > 1 {
			last := payload[len(payload)-1]
			out = append(out, fmt.Sprintf("%v", last[args[1]]))
		}
		fmt.Fprintln(w, strings.Join(out, ","))
	default:
		panic(http.StatusMethodNotAllowed)
	}
}
