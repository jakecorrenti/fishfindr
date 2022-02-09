package db

import (
	"database/sql"
	"errors"

	"github.com/jakecorrenti/fishfindr/types"
	"github.com/mattn/go-sqlite3"
)

// custom errors
var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row doesn't exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

// abstracts the data access layer into a separate structure responsible for storing
// and retrieving data from the database this struct will interact with the SQLite database
type SQLiteRepository struct {
	db *sql.DB
}

// methods do not depend on SQLite at all. they hide database implementation details and provide
// a simple API to interact with any database
type Repository interface {
	Migrate() error
	Create(location types.Location) (*types.Location, error)
	All() ([]types.Location, error)
	GetByID(id string) (*types.Location, error)
	Update(id string, updated types.Location) (*types.Location, error)
	Delete(id string) error
}

// NewSQLiteRepository requires an instance of `sql.DB` type as a dependency. `sql.DB` is an object representing a pool
// of DB connections for all drivers compatible with the `database/sql` interface
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

// Migrate is responsible for migrating the repository. in this case, Migration is creating a
// SQL table and initializing all the data necessary to operate on the repository
//
// this function should be called first before reading or writing data through the repository
func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS locations(
        id TEXT PRIMARY KEY,
        longitude FLOAT NOT NULL,
        latitude FLOAT NOT NULL,
        timestamp TEXT NOT NULL
    )
    `
	// executes a query without returning any rows
	_, err := r.db.Exec(query)

	return err
}

// Create is responsible for taking a row to create, and returns the row after
// insertion or an error if the operation fails
func (r *SQLiteRepository) Create(location types.Location) (*types.Location, error) {
	_, err := r.db.Exec(
		"INSERT INTO locations(id, longitude, latitude, timestamp) values(?,?,?,?)",
		location.Id, location.Longitude, location.Latitude, location.Timestamp,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		// `errors.As` finds the first error in err's chain that matches the target
		// if there is a match, sets target to that error value and returns true,
		// else returns false
		if errors.As(err, &sqliteErr) {
			// `errors.Is` reports whether any error in err's chain matches the target
			// an error is considered to match a target if it is equal to that target
			// or if it implements a method `Is(error) bool` such that `Is(target`
			// returns true
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &location, nil
}

// All is responsible for returning all available records in the Locations repository
func (r *SQLiteRepository) All() ([]types.Location, error) {
	// `Query` executes a query that returns rows, typically a SELECT.
	// the result of `db.Query` is a `sql.Rows` struct that represents a cursor
	// to SQL rows
	rows, err := r.db.Query("SELECT * FROM locations")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var loggedLocations []types.Location
	// `rows.Next()` prepares the enxt result row for reading with the Scan method
	// returns true on success, or false if there is no next result or an error
	// happened while preparing it
	//
	// the error should be consulted to determine between the two possible cases
	for rows.Next() {
		var location types.Location
		// `rows.Scan()` copies the columns in the current row into the values pointed at by dest.
		// the number of values in dest must be the same as the number of columns in Rows
		if err := rows.Scan(&location.Id, &location.Latitude, &location.Longitude, &location.Timestamp); err != nil {
			return nil, err
		}

		loggedLocations = append(loggedLocations, location)
	}

	return loggedLocations, nil
}

// GetByID is responsible for returning the location that is associated with the given Id
func (r *SQLiteRepository) GetByID(id string) (*types.Location, error) {
	// `db.QueryRow` executes a query that returns at most one row. `db.QueryRow` always returns a non-nil value.
	// errors are deferred until `row.Scan` is called
	row := r.db.QueryRow("SELECT * FROM locations WHERE id = ?", id)

	var location types.Location
	if err := row.Scan(&location.Id, &location.Latitude, &location.Longitude, &location.Timestamp); err != nil {
		// `sql.ErrNoRows` indicates that the `db.QueryRow` execution resulted in no rows that match the given id
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &location, nil
}

// Update is responsible for replacing values for a record that match a given id
func (r *SQLiteRepository) Update(id string, updated types.Location) (*types.Location, error) {
	// if id == 0 {
	//     return nil, errors.New("invalid updated ID")
	// }

	res, err := r.db.Exec(
		"UPDATE locations SET id = ?, longitude = ?, latitude = ?, timestamp = ?",
		id, updated.Longitude, updated.Latitude, updated.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	// `res.RowsAffected()` returns the number of rows affected by an update, insert, or delete.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	// if no rows were affected, this is considered a failure
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

// Delete is responsible for deleting the row with the specified Id
func (r *SQLiteRepository) Delete(id string) error {
	res, err := r.db.Exec("DELETE FROM locations WHERE id = ?", id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// if rows affected is none, considered a failure
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
