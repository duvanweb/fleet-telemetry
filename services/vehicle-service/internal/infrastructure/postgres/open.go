package postgres

import "database/sql"

// openDB is a thin wrapper around sql.Open to allow mocking in tests.
func openDB(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}
