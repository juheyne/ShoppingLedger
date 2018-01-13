package model

import (
	"database/sql"
	"time"
	"fmt"
)

type Expense struct {
	ID int `json:"id"`
	Payer string `json:"payer"`
	Amount float32 `json:"amount"`
	Note string `json:"note"`
	Date time.Time `json:"date"`
}

func (e *Expense) Get(db *sql.DB) error {
	err := db.QueryRow("SELECT payer, amount, note, date FROM expenses WHERE id=?", e.ID).Scan(&e.Payer, &e.Amount, &e.Note, &e.Date)

	return err
}

func (e *Expense) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE expenses SET payer=?, amount=?, note=?, date=? WHERE id=?", e.Payer, e.Amount, e.Note, e.Date, e.ID)
	return err
}

func (e *Expense) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM expenses WHERE id=?", e.ID)
	return err
}

func (e *Expense) Create(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO expenses(payer, amount, note, date) VALUES (?,?,?,?)", e.Payer, e.Amount, e.Note, e.Date)
	if err != nil {
		fmt.Println(err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return err
	}

	e.ID = int(id)

	return nil
}

func GetExpenses(db *sql.DB, start, count int) ([]Expense, error) {
	rows, err := db.Query("SELECT id, payer, amount, note, date FROM expenses LIMIT ? OFFSET ?", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	expenses := []Expense{}

	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.Payer, &e.Amount, &e.Note, &e.Date); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}