package sqlite

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"test_backend/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO init db

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(UrlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(UrlToSave, alias)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: url.alias" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(id int64) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Query(id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer res.Close()

	if !res.Next() {
		return "", storage.ErrURLNotFound
	}

	var url string
	if err := res.Scan(&url); err != nil {
		return "", fmt.Errorf("%s: failed to scan result: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(id int64) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
