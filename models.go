package main

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

type Feed struct {
	ID        uuid.UUID
	URL       string
	Channel   string
	Network   string
	CreatedAt time.Time
	DeletedAt NullTime
	LastFetch time.Time

	exists bool
}

type Item struct {
	ID        uuid.UUID
	URL       string
	FeedID    uuid.UUID
	CreatedAt time.Time

	exists bool
}

type BlackList struct {
	ID        uuid.UUID
	FeedID    uuid.UUID
	CreatedAt time.Time

	exists bool
}

func (m *Feed) Create(db DB) error {
	const query = `INSERT INTO feeds (id, url, channel, network) VALUES ($1, $2, $3, $4)`

	m.ID = uuid.New()
	if _, err := db.Exec(query, m.ID, m.URL, m.Channel, m.Network); err != nil {
		return errors.Wrap(err, "failed to create feed")
	}

	return nil
}

func (m *Feed) Delete(db DB) error {
	const query = `DELETE FROM feeds WHERE id = $1`

	if _, err := db.Exec(query, m.ID); err != nil {
		return errors.Wrap(err, "failed to delete feed")
	}

	return nil
}

func FindAllFeeds(db DB) ([]*Feed, error) {
	const query = `SELECT id, url, channel, network, created_at, deleted_at, last_fetch FROM feeds WHERE deleted_at IS NULL`

	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch feeds")
	}

	var feeds []*Feed

	defer rows.Close()

	for rows.Next() {
		var feed Feed

		if err := rows.Scan(&feed.ID, &feed.URL, &feed.Channel, &feed.Network, &feed.CreatedAt, &feed, DeletedAt, &feed.LastFetch); err != nil {
			return nil, errors.Wrap(err, "failed to scan results")
		}

		feeds = append(feeds, &feed)
	}

	return feeds, nil
}
