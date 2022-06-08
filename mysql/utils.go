package mysql

import "database/sql"

// SELECT
func Read(queryStr string, args ...any) (rows *sql.Rows, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	return stmtOut.Query(args...)
}

// INSERT, UPDATE, DELETE
func Write(queryStr string, args ...any) (result sql.Result, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	return stmtIns.Exec(args...)
}

// SELECT
func TxRead(tx *sql.Tx, queryStr string, args ...any) (rows *sql.Rows, err error) {
	stmtOut, err := tx.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	return stmtOut.Query(args...)
}

// INSERT, UPDATE, DELETE
func TxWrite(tx *sql.Tx, queryStr string, args ...any) (result sql.Result, err error) {
	stmtIns, err := tx.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	return stmtIns.Exec(args...)
}
