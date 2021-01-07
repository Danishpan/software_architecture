package postgreSql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"softarch/pkg/models"
	"strconv"
	"time"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	//Query
	s := &models.Snippet{}
	stml := "SELECT id, title, content, created, expires FROM snippets WHERE expires > CLOCK_TIMESTAMP() AND id=$1"

	err := m.DB.QueryRow(context.Background(), stml, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stml := "SELECT id, title, content, created, expires FROM snippets " +
		"WHERE expires > CLOCK_TIMESTAMP() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(context.Background(), stml)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, err
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stml := "INSERT INTO snippets (title, content, created, expires)" + "VALUES($1, $2, $3, $4) RETURNING id"
	d, err := strconv.Atoi(expires)
	if err != nil {
		return 0, err
	}
	var lastId int
	err = m.DB.QueryRow(context.Background(), stml, title, content, time.Now(), time.Now().AddDate(0, 0, d)).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return int(lastId), err
}
