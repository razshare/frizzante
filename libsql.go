package frizzante

import (
	"database/sql"
	"errors"
	"log"
)

func sqlFindNextFallback(dest ...any) bool { return false }
func sqlFindCloseFallback()                {}

type SqlDialect int64

const (
	SqlDialectMysql      SqlDialect = 0
	SqlDialectPostgresql SqlDialect = 1
)

type Sql struct {
	database *sql.DB
	dialect  SqlDialect
	notifier *Notifier
}

// NewSql creates a sql wrapper.
func NewSql() *Sql {
	return &Sql{
		dialect: SqlDialectMysql,
	}
}

// WithNotifier sets the sql notifier.
func (sql *Sql) WithNotifier(notifier *Notifier) {
	sql.notifier = notifier
}

// WithDatabase sets the sql database.
func (sql *Sql) WithDatabase(database *sql.DB) {
	sql.database = database
}

// WithDialect sets the sql dialect.
func (sql *Sql) WithDialect(dialect SqlDialect) {
	sql.dialect = dialect
}

// Execute executes sql queries that don't return rows, typically INSERT, UPDATE, DELETE queries.
func (sql *Sql) Execute(query string, props ...any) *sql.Result {
	transaction, transactionError := sql.database.Begin()
	if transactionError != nil {
		sql.notifier.SendError(transactionError)
		return nil
	}

	statement, statementError := transaction.Prepare(query)
	if nil != statementError {
		sql.notifier.SendError(statementError)
		return nil
	}

	result, execError := statement.Exec(props...)
	if execError != nil {
		sql.notifier.SendError(execError)
		rollbackError := transaction.Rollback()
		if rollbackError != nil {
			sql.notifier.SendError(rollbackError)
		}
		return nil
	}

	commitError := transaction.Commit()
	if commitError != nil {
		sql.notifier.SendError(commitError)
		return nil
	}

	return &result
}

// Find executes a sql query that returns rows, typically a SELECT query.
//
// It returns a next function and a close function.
//
// Use next to project the next row onto dest.
//
// Next will return false if where are no more rows available.
//
// Use close to close the database context and prevent any subsequent enumerations.
//
// Whenever next returns false, the database context is closed automatically as if calling close.
func (sql *Sql) Find(query string, props ...any) (next func(dest ...any) bool, close func()) {
	next = sqlFindNextFallback
	close = sqlFindCloseFallback

	if nil == sql.notifier {
		log.Fatal("sql notifier not provided")
	}

	if nil == sql.database {
		sql.notifier.SendError(errors.New("sql database not provided"))
		return
	}

	statement, statementError := sql.database.Prepare(query)
	if nil != statementError {
		sql.notifier.SendError(statementError)
		return
	}
	defer statement.Close()

	rows, queryError := statement.Query(props...)
	if queryError != nil {
		sql.notifier.SendError(queryError)
		return
	}

	next = func(dest ...any) bool {
		if !rows.Next() {
			return false
		}

		scanError := rows.Scan(dest...)
		if scanError != nil {
			sql.notifier.SendError(scanError)
			return false
		}
		return true
	}
	close = func() {
		err := rows.Close()
		if err != nil {
			sql.notifier.SendError(err)
		}
	}
	return
}
