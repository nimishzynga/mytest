package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
)

type DbConn struct {
	c *sql.DB
}

func NewDbConn(database, user, password string) *DbConn {
	con, err := sql.Open("mymysql", database+"/"+user+"/"+password)
	if err != nil {
		panic(err)
		return nil
	}
	return &DbConn{con}
}

func (d *DbConn) queryDatabase(query string) *sql.Rows {
	rows, err := d.c.Query(query)
	if err != nil {
		fmt.Println("Error in query")
	}
	return rows
}
