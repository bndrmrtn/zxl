package sqlmodule

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bndrmrtn/zxl/lang"
)

var allowedDrivers = []string{"sqlite3", "mysql", "postgres"}

type DB struct {
	lang.Base
	driver string
	dsn    string
	db     *sql.DB
}

func NewDB(driver, dsn string) (*DB, error) {
	// Validate driver
	isValidDriver := false
	for _, d := range allowedDrivers {
		if d == driver {
			isValidDriver = true
			break
		}
	}

	if !isValidDriver {
		return nil, fmt.Errorf("unsupported database driver: %s. Supported drivers: %s",
			driver, strings.Join(allowedDrivers, ", "))
	}

	// Configure DSN based on driver
	configuredDSN := configureDSN(driver, dsn)

	// Open database connection
	db, err := sql.Open(driver, configuredDSN)
	if err != nil {
		return nil, err
	}

	// Apply driver-specific configurations
	configureDBConnection(db, driver)

	// Verify connection works
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{
		Base:   lang.NewBase("conn", nil),
		driver: driver,
		dsn:    configuredDSN,
		db:     db,
	}, nil
}

// configureDSN adjusts DSN for specific drivers
func configureDSN(driver, dsn string) string {
	switch driver {
	case "mysql":
		// Format: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		if !strings.Contains(dsn, "parseTime=") {
			if strings.Contains(dsn, "?") {
				dsn += "&parseTime=true"
			} else {
				dsn += "?parseTime=true"
			}
		}

		if !strings.Contains(dsn, "charset=") {
			if strings.Contains(dsn, "?") {
				dsn += "&charset=utf8mb4"
			} else {
				dsn += "?charset=utf8mb4"
			}
		}

		// Add timeout parameters
		if !strings.Contains(dsn, "timeout=") {
			if strings.Contains(dsn, "?") {
				dsn += "&timeout=30s"
			} else {
				dsn += "?timeout=30s"
			}
		}

		if !strings.Contains(dsn, "readTimeout=") {
			if strings.Contains(dsn, "?") {
				dsn += "&readTimeout=30s"
			} else {
				dsn += "?readTimeout=30s"
			}
		}

		if !strings.Contains(dsn, "writeTimeout=") {
			if strings.Contains(dsn, "?") {
				dsn += "&writeTimeout=30s"
			} else {
				dsn += "?writeTimeout=30s"
			}
		}

	case "postgres":
		// Format: postgres://username:password@localhost:5432/database_name?sslmode=disable
		if !strings.Contains(dsn, "sslmode=") {
			if strings.Contains(dsn, "?") {
				dsn += "&sslmode=prefer"
			} else {
				dsn += "?sslmode=prefer"
			}
		}

		// Add connection timeout
		if !strings.Contains(dsn, "connect_timeout=") {
			if strings.Contains(dsn, "?") {
				dsn += "&connect_timeout=10"
			} else {
				dsn += "?connect_timeout=10"
			}
		}

	case "sqlite3":
		// Format: file:test.db?cache=shared&mode=memory
		if !strings.Contains(dsn, "cache=") {
			if strings.Contains(dsn, "?") {
				dsn += "&cache=shared"
			} else {
				dsn += "?cache=shared"
			}
		}

		if !strings.Contains(dsn, "_journal=") {
			if strings.Contains(dsn, "?") {
				dsn += "&_journal=WAL"
			} else {
				dsn += "?_journal=WAL"
			}
		}
	}

	return dsn
}

// configureDBConnection sets connection pool settings and other driver-specific options
func configureDBConnection(db *sql.DB, driver string) {
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	// Driver-specific initial queries can be done here if needed
	if driver == "sqlite3" {
		db.Exec("PRAGMA foreign_keys = ON")
		db.Exec("PRAGMA busy_timeout = 5000")
	}
}

func (db *DB) Type() lang.ObjType {
	return lang.TInstance
}

func (*DB) TypeString() string {
	return "sql.db"
}

func (db *DB) Value() any {
	return db
}

func (db *DB) Method(name string) lang.Method {
	switch name {
	case "query":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[1].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}
			rows, err := db.db.Query(args[0].Value().(string), queryArgs...)
			if err != nil {
				fmt.Println("err", queryArgs)
				return nil, err
			}
			defer rows.Close()

			var result []lang.Object
			cols, err := rows.Columns()
			if err != nil {
				return nil, err
			}
			keys := make([]lang.Object, len(cols))
			for i, colName := range cols {
				keys[i] = lang.NewString("column", colName, nil)
			}
			for rows.Next() {
				rowValues := make([]interface{}, len(cols))
				pointers := make([]interface{}, len(cols))
				for i := range rowValues {
					pointers[i] = &rowValues[i]
				}
				if err := rows.Scan(pointers...); err != nil {
					return nil, err
				}
				values := make([]lang.Object, len(cols))
				for i := range rowValues {
					value, err := lang.FromValue(rowValues[i])
					if err != nil {
						return nil, err
					}
					values[i] = value
				}
				result = append(result, lang.NewArray("row", nil, keys, values))
			}
			if err = rows.Err(); err != nil {
				return nil, err
			}

			return lang.NewList("result", result, nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString}).WithVariadicArg("values")

	case "exec":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[1].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			result, err := db.db.Exec(args[0].Value().(string), queryArgs...)
			if err != nil {
				return nil, err
			}

			lastID, _ := result.LastInsertId()
			rowsAffected, _ := result.RowsAffected()

			return lang.NewArray("execResult", nil, []lang.Object{
				lang.NewString("lastInsertId", "lastInsertId", nil),
				lang.NewString("rowsAffected", "rowsAffected", nil),
			}, []lang.Object{
				lang.NewInteger("lastInsertId", int(lastID), nil),
				lang.NewInteger("rowsAffected", int(rowsAffected), nil),
			}), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString}).WithVariadicArg("values")

	case "queryRow":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[1].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			// First get column names through a separate query
			stmt, err := db.db.Prepare(args[0].Value().(string))
			if err != nil {
				return nil, err
			}
			defer stmt.Close()

			rows, err := stmt.Query(queryArgs...)
			if err != nil {
				return nil, err
			}

			cols, err := rows.Columns()
			rows.Close()
			if err != nil {
				return nil, err
			}

			// Now execute the actual queryRow
			row := db.db.QueryRow(args[0].Value().(string), queryArgs...)

			rowValues := make([]interface{}, len(cols))
			pointers := make([]interface{}, len(cols))
			for i := range rowValues {
				pointers[i] = &rowValues[i]
			}

			if err := row.Scan(pointers...); err != nil {
				if err == sql.ErrNoRows {
					return lang.NewList("emptyList", nil, nil), nil
				}
				return nil, err
			}

			keys := make([]lang.Object, len(cols))
			values := make([]lang.Object, len(cols))

			for i, colName := range cols {
				keys[i] = lang.NewString("column", colName, nil)
				value, err := lang.FromValue(rowValues[i])
				if err != nil {
					return nil, err
				}
				values[i] = value
			}

			return lang.NewArray("row", nil, keys, values), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString}).WithVariadicArg("values")

	case "prepare":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			stmt, err := db.db.Prepare(args[0].Value().(string))
			if err != nil {
				return nil, err
			}

			return &Statement{
				Base:  lang.NewBase("stmt", nil),
				stmt:  stmt,
				query: args[0].Value().(string),
			}, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString})

	case "beginTx":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			tx, err := db.db.Begin()
			if err != nil {
				return nil, err
			}

			return &Transaction{
				Base: lang.NewBase("tx", nil),
				tx:   tx,
			}, nil
		})

	case "close":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			if err := db.db.Close(); err != nil {
				return nil, err
			}
			return lang.NewBool("closed", true, nil), nil
		})

	case "ping":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			err := db.db.Ping()
			return lang.NewBool("ping", err == nil, nil), nil
		})

	default:
		return nil
	}
}

func (db *DB) Methods() []string {
	return []string{"close", "query", "queryRow", "exec", "prepare", "beginTx", "ping"}
}

func (db *DB) Variable(variable string) lang.Object {
	switch variable {
	default:
		return nil
	case "$addr":
		return lang.Addr(db)
	case "$driver":
		return lang.NewString("driver", db.driver, nil)
	case "$dsn":
		return lang.NewString("dsn", db.dsn, nil)
	}
}

func (db *DB) Variables() []string {
	return []string{"$addr", "$driver", "$dsn"}
}

func (db *DB) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (db *DB) String() string {
	return fmt.Sprintf("<SQL %s %s>", db.driver, lang.Addr(db))
}

func (db *DB) Copy() lang.Object {
	return db
}

// Statement represents a prepared statement
type Statement struct {
	lang.Base
	stmt  *sql.Stmt
	query string
}

func (s *Statement) Type() lang.ObjType {
	return lang.TInstance
}

func (*Statement) TypeString() string {
	return "sql.statement"
}

func (s *Statement) Value() any {
	return s
}

func (s *Statement) Method(name string) lang.Method {
	switch name {
	case "exec":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[0].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			result, err := s.stmt.Exec(queryArgs...)
			if err != nil {
				return nil, err
			}

			lastID, _ := result.LastInsertId()
			rowsAffected, _ := result.RowsAffected()

			return lang.NewArray("execResult", nil, []lang.Object{
				lang.NewString("lastInsertId", "lastInsertId", nil),
				lang.NewString("rowsAffected", "rowsAffected", nil),
			}, []lang.Object{
				lang.NewInteger("lastInsertId", int(lastID), nil),
				lang.NewInteger("rowsAffected", int(rowsAffected), nil),
			}), nil
		}).WithVariadicArg("values")

	case "query":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[0].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			rows, err := s.stmt.Query(queryArgs...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var result []lang.Object
			cols, err := rows.Columns()
			if err != nil {
				return nil, err
			}
			keys := make([]lang.Object, len(cols))
			for i, colName := range cols {
				keys[i] = lang.NewString("column", colName, nil)
			}
			for rows.Next() {
				rowValues := make([]interface{}, len(cols))
				pointers := make([]interface{}, len(cols))
				for i := range rowValues {
					pointers[i] = &rowValues[i]
				}
				if err := rows.Scan(pointers...); err != nil {
					return nil, err
				}
				values := make([]lang.Object, len(cols))
				for i := range rowValues {
					value, err := lang.FromValue(rowValues[i])
					if err != nil {
						return nil, err
					}
					values[i] = value
				}
				result = append(result, lang.NewArray("row", nil, keys, values))
			}
			if err = rows.Err(); err != nil {
				return nil, err
			}

			return lang.NewList("result", result, nil), nil
		}).WithVariadicArg("values")

	case "queryRow":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[0].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			// Get column names
			rows, err := s.stmt.Query(queryArgs...)
			if err != nil {
				return nil, err
			}

			cols, err := rows.Columns()
			rows.Close()
			if err != nil {
				return nil, err
			}

			// Execute the actual queryRow
			row := s.stmt.QueryRow(queryArgs...)

			rowValues := make([]interface{}, len(cols))
			pointers := make([]interface{}, len(cols))
			for i := range rowValues {
				pointers[i] = &rowValues[i]
			}

			if err := row.Scan(pointers...); err != nil {
				if err == sql.ErrNoRows {
					return lang.NewList("emptyList", nil, nil), nil
				}
				return nil, err
			}

			keys := make([]lang.Object, len(cols))
			values := make([]lang.Object, len(cols))

			for i, colName := range cols {
				keys[i] = lang.NewString("column", colName, nil)
				value, err := lang.FromValue(rowValues[i])
				if err != nil {
					return nil, err
				}
				values[i] = value
			}

			return lang.NewArray("row", nil, keys, values), nil
		}).WithVariadicArg("values")

	case "close":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			if err := s.stmt.Close(); err != nil {
				return nil, err
			}
			return lang.NewBool("closed", true, nil), nil
		})

	default:
		return nil
	}
}

func (s *Statement) Methods() []string {
	return []string{"exec", "query", "queryRow", "close"}
}

func (s *Statement) Variable(variable string) lang.Object {
	switch variable {
	default:
		return nil
	case "$addr":
		return lang.Addr(s)
	case "$query":
		return lang.NewString("query", s.query, nil)
	}
}

func (s *Statement) Variables() []string {
	return []string{"$addr", "$query"}
}

func (s *Statement) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (s *Statement) String() string {
	return fmt.Sprintf("<SQL Statement %s>", lang.Addr(s))
}

func (s *Statement) Copy() lang.Object {
	return s
}

// Transaction represents a database transaction
type Transaction struct {
	lang.Base
	tx *sql.Tx
}

func (t *Transaction) Type() lang.ObjType {
	return lang.TInstance
}

func (*Transaction) TypeString() string {
	return "sql.transaction"
}

func (t *Transaction) Value() any {
	return t
}

func (t *Transaction) Method(name string) lang.Method {
	switch name {
	case "commit":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			if err := t.tx.Commit(); err != nil {
				return nil, err
			}
			return lang.NewBool("committed", true, nil), nil
		})

	case "rollback":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			if err := t.tx.Rollback(); err != nil {
				return nil, err
			}
			return lang.NewBool("rolledback", true, nil), nil
		})

	case "exec":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[1].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			result, err := t.tx.Exec(args[0].Value().(string), queryArgs...)
			if err != nil {
				return nil, err
			}

			lastID, _ := result.LastInsertId()
			rowsAffected, _ := result.RowsAffected()

			return lang.NewArray("execResult", nil, []lang.Object{
				lang.NewString("lastInsertId", "lastInsertId", nil),
				lang.NewString("rowsAffected", "rowsAffected", nil),
			}, []lang.Object{
				lang.NewInteger("lastInsertId", int(lastID), nil),
				lang.NewInteger("rowsAffected", int(rowsAffected), nil),
			}), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString}).WithVariadicArg("values")

	case "query":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			var queryArgs []interface{}
			for _, arg := range args[1].Value().([]lang.Object) {
				queryArgs = append(queryArgs, arg.Value())
			}

			rows, err := t.tx.Query(args[0].Value().(string), queryArgs...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var result []lang.Object
			cols, err := rows.Columns()
			if err != nil {
				return nil, err
			}
			keys := make([]lang.Object, len(cols))
			for i, colName := range cols {
				keys[i] = lang.NewString("column", colName, nil)
			}
			for rows.Next() {
				rowValues := make([]interface{}, len(cols))
				pointers := make([]interface{}, len(cols))
				for i := range rowValues {
					pointers[i] = &rowValues[i]
				}
				if err := rows.Scan(pointers...); err != nil {
					return nil, err
				}
				values := make([]lang.Object, len(cols))
				for i := range rowValues {
					value, err := lang.FromValue(rowValues[i])
					if err != nil {
						return nil, err
					}
					values[i] = value
				}
				result = append(result, lang.NewArray("row", nil, keys, values))
			}
			if err = rows.Err(); err != nil {
				return nil, err
			}

			return lang.NewList("result", result, nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString}).WithVariadicArg("values")

	case "prepare":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			stmt, err := t.tx.Prepare(args[0].Value().(string))
			if err != nil {
				return nil, err
			}

			return &Statement{
				Base:  lang.NewBase("stmt", nil),
				stmt:  stmt,
				query: args[0].Value().(string),
			}, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString})

	default:
		return nil
	}
}

func (t *Transaction) Methods() []string {
	return []string{"commit", "rollback", "exec", "query", "prepare"}
}

func (t *Transaction) Variable(variable string) lang.Object {
	switch variable {
	default:
		return nil
	case "$addr":
		return lang.Addr(t)
	}
}

func (t *Transaction) Variables() []string {
	return []string{"$addr"}
}

func (t *Transaction) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (t *Transaction) String() string {
	return fmt.Sprintf("<SQL Transaction %s>", lang.Addr(t))
}

func (t *Transaction) Copy() lang.Object {
	return t
}
