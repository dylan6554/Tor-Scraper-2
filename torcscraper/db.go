package main

import "database/sql"

func insertRecord(db *sql.DB, r Record) error {
	_, err := db.Exec(`
		INSERT INTO findings
		(url, title, source_title, category, criticality, content, screenshot)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`,
		r.URL,
		r.Title,
		r.SourceTitle,
		r.Category,
		r.Criticality,
		r.Content,
		r.Screenshot,
	)
	return err
}
