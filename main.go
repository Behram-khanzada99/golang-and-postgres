package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

func insertRecord(db *sql.DB, name string) (int, error) {
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
	// Check if the record with the specified ID exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM employees WHERE id = $1", id).Scan(&count)
	if err != nil {
		return err
	}

	// If the record doesn't exist, return an error
	if count == 0 {
		return fmt.Errorf("record with ID %d does not exist", id)
	}

	// Update the record if it exists
	_, err = db.Exec("UPDATE employees SET name = $1 WHERE id = $2", newName, id)
	return err
}

func deleteRecord(db *sql.DB, id int) error {
	// Check if the record with the specified ID exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM employees WHERE id = $1", id).Scan(&count)
	if err != nil {
		return err
	}

	// If the record doesn't exist, return an error
	if count == 0 {
		return fmt.Errorf("record with ID %d does not exist", id)
	}

	// Delete the record if it exists
	_, err = db.Exec("DELETE FROM employees WHERE id = $1", id)
	return err
}

// inserted timestamp into create table function

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id serial PRIMARY KEY,
			name VARCHAR (255) NOT NULL,
			created_at TIMESTAMP DEFAULT current_timestamp
		)
	`)
	return err
}

func insertAndReadRecords(db *sql.DB) error {
	startTime := time.Now()

	// Insert records
	for i := 0; i < 100000; i++ {
		name := fmt.Sprintf("User%d", i)
		_, err := insertRecord(db, name)
		if err != nil {
			return err
		}
	}

	insertElapsedTime := time.Since(startTime)
	fmt.Printf("Inserted 100,000 records in %v\n", insertElapsedTime)

	// Read records
	startTime = time.Now()

	records, err := getAllRecords(db)
	if err != nil {
		return err
	}

	readElapsedTime := time.Since(startTime)
	fmt.Printf("Read %d records in %v\n", len(records), readElapsedTime)

	return nil
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

	// Create the "employees" table if it doesn't exist
	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'employees' created or already exists")

	// Perform CRUD operations
	var choice int

	for {
		fmt.Println("\nChoose an operation:")
		fmt.Println("Press 1 to Insert Record")
		fmt.Println("Press 2 to Update Record")
		fmt.Println("Press 3 to Delete Record")
		fmt.Println("Press 4 to Read All Records")
		fmt.Println("Press 5 to Insert and Read 100k Records")
		// fmt.Println("Press 6 to insert and Read 100k Records Concurrently")
		fmt.Println("Press 6 to Exit")

		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)
		// fmt.Scanf("%s", &choice)

		switch choice {
		case 1:
			newID, err := insertRecord(db, "John")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created a new record with ID: %d\n", newID)

		case 2:
			var id int
			var newName string

			fmt.Print("Enter the ID to update: ")
			fmt.Scan(&id)

			fmt.Print("Enter the new name: ")
			fmt.Scan(&newName)

			err := updateRecord(db, id, newName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Record updated successfully")

		case 3:
			var idToDelete int

			fmt.Print("Enter the ID to delete: ")
			fmt.Scan(&idToDelete)

			err := deleteRecord(db, idToDelete)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Record deleted successfully")

		case 4:
			records, err := getAllRecords(db)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("All Records:")
			for _, record := range records {
				fmt.Printf("ID: %d, Name: %s\n", record.ID, record.Name)
			}

		case 5:
			// Call the insertAndReadRecords function
			err := insertAndReadRecords(db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted and Read Records successfully")

		case 6:
			fmt.Println("Exiting program.")
			return

		default:
			fmt.Println("Invalid choice. Please enter a valid option.")
		}

	}
}
