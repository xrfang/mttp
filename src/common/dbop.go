package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func RangeRows(rows *sql.Rows, proc func()) {
	defer func() {
		if e := recover(); e != nil {
			rows.Close()
			panic(e)
		}
	}()
	for rows.Next() {
		proc()
	}
	assert(rows.Err())
}

func FetchRows(rows *sql.Rows) []map[string]interface{} {
	defer func() {
		if e := recover(); e != nil {
			rows.Close()
			panic(e)
		}
	}()
	cols, err := rows.Columns()
	assert(err)
	raw := make([][]byte, len(cols))
	ptr := make([]interface{}, len(cols))
	for i := range raw {
		ptr[i] = &raw[i]
	}
	var recs []map[string]interface{}
	for rows.Next() {
		assert(rows.Scan(ptr...))
		rec := make(map[string]interface{})
		for i, r := range raw {
			if r == nil {
				rec[cols[i]] = nil
			} else {
				rec[cols[i]] = string(r)
			}
		}
		recs = append(recs, rec)
	}
	assert(rows.Err())
	return recs
}
