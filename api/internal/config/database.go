package config

import (
	"errors"
	"fmt"
	"net/url"
)

type DB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`

	TLS bool `json:"tls"`
}

func (db *DB) FromURI(uri string) error {
	connect, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("failed to parse database URI: %w", err)
	}

	db.Host = connect.Hostname()
	db.Port = 5432

	if connect.Port() != "" {
		db.Port = 5432
	}

	if connect.User != nil {
		db.User = connect.User.Username()
		db.Password, _ = connect.User.Password()
	}

	if connect.Path != "" {
		db.Database = connect.Path[1:]
	}

	if connect.Query().Get("sslmode") == "require" {
		db.TLS = true
	}

	return nil
}

func (db *DB) URI() string {
	connect := url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", db.Host, db.Port),
	}

	if db.User != "" {
		connect.User = url.UserPassword(db.User, db.Password)
	}

	if db.Database != "" {
		connect.Path = db.Database
	}

	if db.TLS {
		connect.RawQuery = url.Values{
			"sslmode": []string{"require"},
		}.Encode()
	} else {
		connect.RawQuery = url.Values{
			"sslmode": []string{"disable"},
		}.Encode()
	}

	return connect.String()
}

func (db *DB) Validate() error {
	if db.Host == "" {
		return errors.New("database host is required")
	}

	if db.Port == 0 {
		return errors.New("database port is required")
	}

	return nil
}
