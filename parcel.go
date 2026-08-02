package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func createSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER NOT NULL,
			status TEXT NOT NULL,
			address TEXT NOT NULL,
			created_at TEXT NOT NULL
		)
	`)

	return err
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		`INSERT INTO parcel (client, status, address, created_at)
		VALUES (?, ?, ?, ?)`,
		p.Client, p.Status, p.Address, p.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow(
		`SELECT number, client, status, address, created_at
		FROM parcel
		WHERE number = ?`,
		number,
	)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(
		`SELECT number, client, status, address, created_at
		FROM parcel
		WHERE client = ?`,
		client,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(
		`UPDATE parcel
		SET status = ?
		WHERE number = ?`,
		status, number,
	)

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	if p.Status != ParcelStatusRegistered {
		return errors.New("address can be changed only for registered parcel")
	}

	_, err = s.db.Exec(
		`UPDATE parcel
		SET address = ?
		WHERE number = ?`,
		address, number,
	)

	return err
}

func (s ParcelStore) Delete(number int) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	if p.Status != ParcelStatusRegistered {
		return errors.New("only registered parcel can be deleted")
	}

	_, err = s.db.Exec(
		`DELETE FROM parcel
		WHERE number = ?`,
		number,
	)

	return err
}
