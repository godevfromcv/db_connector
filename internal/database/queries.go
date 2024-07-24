package database

import (
	"database/sql"
	"fmt"
)

func ExecuteQuery(db *sql.DB, query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func PrintQueryResults(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(values))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return err
		}

		for i, colName := range columns {
			fmt.Printf("%s: %v\n", colName, values[i])
		}
		fmt.Println("-------------------------")
	}

	return rows.Err()
}

func GetTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}
