package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var mu sync.Mutex // Mutex for synchronization

func createDBConnection() (*sql.DB, error) {
	// Replace with your PostgreSQL connection string
	dbURL := "postgres://postgres:alpha123@localhost:5432/my_pgdb?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

func insertRecord(db *sql.DB, name string) (int, error) {
	mu.Lock()         // Lock before starting the critical section
	defer mu.Unlock() // Ensure the mutex is unlocked even if an error occurs

	var id int
	err := db.QueryRow("INSERT INTO employees(name) VALUES($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getAllRecords(db *sql.DB) ([]Record, error) {
	mu.Lock()
	defer mu.Unlock()

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
	mu.Lock()
	defer mu.Unlock()

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
	mu.Lock()
	defer mu.Unlock()

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

// Insert records concurrently
func insertAndReadRecordsConcurrently(db *sql.DB) error {
	startTime := time.Now()

	var wg sync.WaitGroup
	// var mu sync.Mutex // Mutex for synchronization

	// Insert records concurrently
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			name := fmt.Sprintf("User%d", index)
			_, err := insertRecord(db, name)
			if err != nil {
				log.Println(err)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	insertElapsedTime := time.Since(startTime)
	fmt.Printf("Inserted 100000 records in %v\n", insertElapsedTime)

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
		fmt.Println("1. Insert Record")
		fmt.Println("2. Update Record")
		fmt.Println("3. Delete Record")
		fmt.Println("4. Read All Records")
		fmt.Println("5. Insert and Read 100k Records Concurrently")
		fmt.Println("6. Exit")

		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

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
			// Call the insertAndReadRecordsConcurrently function
			err := insertAndReadRecordsConcurrently(db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted and Read Records concurrently")

		case 6:
			fmt.Println("Exiting program.")
			return

		default:
			fmt.Println("Invalid choice. Please enter a valid option.")
		}
	}
}
