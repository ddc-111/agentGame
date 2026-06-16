//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	db, _ := sql.Open("sqlite", "game.db")
	defer db.Close()
	rows, _ := db.Query("SELECT id, items FROM players WHERE deleted_at IS NULL")
	defer rows.Close()
	for rows.Next() {
		var id int
		var items string
		rows.Scan(&id, &items)
		fmt.Printf("Player %d items: [%s]\n", id, items)
	}
}
