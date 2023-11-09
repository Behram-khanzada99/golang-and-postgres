package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func createDBConnection() (*sql.DB, error) {
	dbURL := "postgres://postgres:alpha123@localhost:5432/my_pgdb?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createRecord(db *sql.DB, name string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO employees(name) VALUES($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getAllRecords(db *sql.DB) ([]Record, error) {
	rows, err := db.Query("SELECT id, name FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

func updateRecord(db *sql.DB, id int, newName string) error {
	_, err := db.Exec("UPDATE employees SET name = $1 WHERE id = $2", newName, id)
	return err
}

func deleteRecord(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
	return err
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id serial PRIMARY KEY,
			name VARCHAR (255) NOT NULL
		)
	`)
	return err
}

type Record struct {
	ID   int
	Name string
}

func main() {
	db, err := createDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL database")

	// Create the "employees" if it doesn't exist
	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'employees' created or already exists")

	//Perform CRUD operations

	newID, err := createRecord(db, "Behram")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created a new record with ID: %d\n", newID)

	records, err := getAllRecords(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		fmt.Printf("ID: %d, Name: %s\n", record.ID, record.Name)
	}

	err = updateRecord(db, 1, "Tomas")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record updated successfully")

	err = deleteRecord(db, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record with ID 2 deleted successfully")
}
