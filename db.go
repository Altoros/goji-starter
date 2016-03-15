package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
)

type DbConnection struct {
	connString string
}

func NewDBConnection(connString string) (conn *DbConnection) {
	conn = new(DbConnection)
	conn.createLocalConnection(connString)
	return
}

func (c *DbConnection) createLocalConnection(connString string) (err error) {
	db, err := sql.Open("postgres", connString)
	defer db.Close()
	if err == nil {
		err = db.Ping()
		if err == nil {
			c.connString = connString
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS comments (id SERIAL, data TEXT)")
		if err != nil {
			log.Printf("Error creating comments table - %s", err)
		}
	}
	if err != nil {
		log.Printf("Connection failed - %s", err)
	}
	return
}

func (c *DbConnection) GetComments() (comments []comment, err error) {
	comments = make([]comment, 0)
	db, err := sql.Open("postgres", c.connString)
	defer db.Close()
	if err != nil {
		log.Printf("Connection to database failed - %s", err)
		return
	}
	queryStmt, err := db.Prepare("SELECT data FROM comments")
	if err != nil {
		return
	}
	rows, err := queryStmt.Query()
	defer rows.Close()
	for rows.Next() {
		var commentData []byte
		var comment comment
		if err := rows.Scan(&commentData); err == nil {
			err := json.Unmarshal(commentData, &comment)
			if err != nil {
				log.Printf("Error unmarshaling comment - %s", err)
			} else {
				comments = append(comments, comment)
			}
		}
	}

	return
}

func (c *DbConnection) AddComment(cmt comment) (err error) {
	db, err := sql.Open("postgres", c.connString)
	defer db.Close()
	if err != nil {
		log.Printf("Connection to database failed - %s", err)
		return
	}
	commentData, err := json.Marshal(cmt)
	if err == nil {
		_, err = db.Exec("INSERT INTO comments(data) VALUES ($1)", commentData)
	}
	return
}
