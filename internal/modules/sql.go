package modules

import "database/sql"

type SQL struct {
	conn *sql.Conn
}

func NewSQL() *SQL {
	return &SQL{}
}

func (*SQL) Namespace() string {
	return "sql"
}
