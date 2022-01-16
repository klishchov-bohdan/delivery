package transaction

import (
	"database/sql"
	"strconv"
	"strings"
)

type PipelineStmt struct {
	query string
	args  []interface{}
}

func NewPipelineStmt(query string, args ...interface{}) *PipelineStmt {
	return &PipelineStmt{query, args}
}

func (ps *PipelineStmt) Exec(tx Transaction, lastInsertId int64) (sql.Result, error) {
	query := strings.Replace(ps.query, "{LAST_INS_ID}", strconv.Itoa(int(lastInsertId)), -1)
	return tx.Exec(query, ps.args...)
}

func RunPipeline(tx Transaction, stmts ...*PipelineStmt) (sql.Result, error) {
	var res sql.Result
	var err error
	var lastInsId int64

	for _, ps := range stmts {
		res, err = ps.Exec(tx, lastInsId)
		if err != nil {
			return nil, err
		}

		lastInsId, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
