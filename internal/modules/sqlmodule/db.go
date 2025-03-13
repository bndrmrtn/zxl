package sqlmodule

import (
	"database/sql"
	"fmt"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type DB struct {
	lang.Base

	driver string
	dsn    string
	db     *sql.DB
}

func NewDB(driver, dsn string) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		Base:   lang.NewBase("conn", nil),
		driver: driver,
		dsn:    dsn,
		db:     db,
	}, nil
}

func (db *DB) Type() lang.ObjType {
	return lang.TDefinition
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
			rows, err := db.db.Query(args[0].Value().(string))
			if err != nil {
				return nil, err
			}

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

			return lang.NewList("result", result, nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "query", Type: lang.TString})
	case "close":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			if err := db.db.Close(); err != nil {
				return nil, err
			}
			return lang.NewBool("closed", true, nil), nil
		})
	default:
		return nil
	}
}

func (db *DB) Methods() []string {
	return []string{"close", "query"}
}

func (db *DB) Variable(variable string) lang.Object {
	switch variable {
	default:
		return nil
	case "$addr":
		return lang.Addr(db)
	}
}

func (db *DB) Variables() []string {
	return []string{"$addr"}
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
