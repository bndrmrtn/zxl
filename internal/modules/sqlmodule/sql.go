package sqlmodule

import (
	"database/sql"
	"fmt"

	"github.com/bndrmrtn/zxl/internal/lang"
	_ "github.com/mattn/go-sqlite3"
)

type SQL struct {
	conn *sql.Conn
}

func New() *SQL {
	return &SQL{}
}

func (*SQL) Namespace() string {
	return "sql"
}

func (s *SQL) Objects() map[string]lang.Object {
	return nil
}

func (s *SQL) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"open": lang.NewFunction([]string{"driver", "dsn"}, func(args []lang.Object) (lang.Object, error) {
			if args[0].Type() != lang.TString || args[1].Type() != lang.TString {
				return nil, fmt.Errorf("invalid argument types")
			}
			return NewDB(args[0].Value().(string), args[1].Value().(string))
		}, nil),
	}
}
