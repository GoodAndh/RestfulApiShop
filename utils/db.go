package utils

import (
	"database/sql"
	"restful/exception"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		rollback:= tx.Rollback()
		exception.PanicIfError(rollback)
		
	} else {
		commit:=tx.Commit()
		exception.PanicIfError(commit)


	}
}
