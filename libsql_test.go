package frizzante

import (
	"testing"
)

func TestSqlOperatorCreate(test *testing.T) {
	sql := NewSql()
	if nil == sql {
		test.Fatal("could not create sql")
	}
}
