// SQLI by db.Query(some).Scan(&other)
package main

import (
	"database/sql"
	"fmt"
	"os"
)

func main() {
	var name string
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	q := fmt.Sprintf("SELECT name FROM users where id = '%s'", os.Args[1])
	row := db.QueryRow(q)
	err = row.Scan(&name)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
