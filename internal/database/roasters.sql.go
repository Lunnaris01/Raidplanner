package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

/*
Roasters:
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    server TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

*/

type Roaster struct {
	RoasterID int
	Username  string
	Server    string
	CreatedAt time.Time
}

func (q *Queries) AddRoaster(ctx context.Context, Username string, Server string) (int, error) {
	// Define the query string
	queryString := "INSERT INTO roasters (username, server) VALUES (?, ?) RETURNING id"

	stmt, err := q.PrepareContext(ctx, queryString)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close() // Ensure the statement is closed after use
	// Execute the statement with the provided parameters
	game, err := stmt.ExecContext(ctx, Username, Server)
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}
	roasterID, err := game.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	log.Println("Roaster added successfully")
	return int(roasterID), nil
}

func (q *Queries) DeleteRoaster(ctx context.Context, ID int) error {
	// Define the query string
	queryString := "DELETE from roasters where id=?"

	stmt, err := q.PrepareContext(ctx, queryString)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close() // Ensure the statement is closed after use
	// Execute the statement with the provided parameters
	_, err = stmt.ExecContext(ctx, ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	log.Println("Roaster deleted successfully")
	return nil
}

func (q *Queries) GetRoasters(ctx context.Context) ([]Roaster, error) {
	// Define the query string
	queryString := "SELECT id, username, server, created_at from roasters"

	stmt, err := q.PrepareContext(ctx, queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close() // Ensure the statement is closed after use
	// Execute the statement with the provided parameters
	var roasters []Roaster
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data from database: %w", err)
	}
	for rows.Next() {
		var i Roaster
		if err := rows.Scan(
			&i.RoasterID,
			&i.Username,
			&i.Server,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		roasters = append(roasters, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roasters, nil
}
