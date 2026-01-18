package main

import "database/sql"

var db *sql.DB

type Record struct {
	URL         string
	Title       string
	SourceTitle string
	Category    string
	Criticality string
	Content     string
	Screenshot  string
}
