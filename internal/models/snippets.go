package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	// for multi-line SQL statements, it's often clearer to use backticks to create a raw string literal
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// .Exec() returns a sql.Result and an error.
	// The sql.Result contains information about the effect of the statement,
	// such as the ID of the last inserted row or the number of rows affected by an UPDATE or DELETE statement.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID is returned as an int64, but our Snippet struct uses an int for the ID field.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
	FROM snippets
	WHERE id = ? AND expires > UTC_TIMESTAMP()`

	// returns a single row, which is represented by a sql.Row object.
	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	// notice that we use the pointer to the Snippet struct here,
	// so that the Scan() method can populate the fields of the struct directly.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if no matching record is found, the Scan() method will return a sql.ErrNoRows error.
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	return s, nil
}

// func (m *SnippetModel) Get(id int) (*Snippet, error) {
// 	stmt := `SELECT id, title, content, created, expires
// 	FROM snippets
// 	WHERE id = ? AND expires > UTC_TIMESTAMP()`

// 	row := m.DB.QueryRow(stmt, id)

// 	s := &Snippet{}

// 	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, ErrNoRecord
// 		}
// 		return nil, err
// 	}

// 	return s, nil
// }
