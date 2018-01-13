package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type expense struct {
	ID int `json:"id"`
	Payer int `json:"payer"`
	Amount float32 `json:"amount"`
	Note string `json:"note"`
}

func (e *expense) getExpense(db *sql.DB) error {
	err := db.QueryRow("SELECT payer, amount, note FROM expenses WHERE id=?", e.ID).Scan(&e.Payer, &e.Amount, &e.Note)

	return err
}

func (e *expense) updateExpense(db *sql.DB) error {
	_, err := db.Exec("UPDATE expenses SET payer=?, amount=?, note=? WHERE id=?", e.Payer, e.Amount, e.Note, e.ID)
	return err
}

func (e *expense) deleteExpense(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM expenses WHERE id=?", e.ID)
	return err
}

func (e *expense) createProduct(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO expenses(payer, amount, note) VALUES (?,?,?)", e.Payer, e.Amount, e.Note)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = int(id)

	return nil
}

func getExpenses(db *sql.DB, start, count int) ([]expense, error) {
	rows, err := db.Query("SELECT id, payer, amount, note FROM expenses LIMIT ? OFFSET ?", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	expenses := []expense{}

	for rows.Next() {
		var e expense
		if err := rows.Scan(&e.ID, &e.Payer, &e.Amount, &e.Note); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}