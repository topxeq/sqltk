package sqltk

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/topxeq/tk"
)

// use SQLite3 for test
// use test.db in current directory as SQLite3 file
var dbConnectString = `test.db`

func Test001(t *testing.T) {
	// connect SQLite3 database
	dbT, errT := ConnectDB("sqlite3", dbConnectString)

	// exit if failed to connect
	tk.CheckErrf("failed to connect: %v", errT)

	defer dbT.Close()

	// ExecV is used to run SQL command without row result
	// such as DROP TABLE, CREATE TABLE, INSERT or UPDATE...
	insertIDT, rowsAffectedT, errT := ExecV(dbT, `DROP TABLE TXTEST`)

	// tk.CheckErrf("failed to drop DB: %v", errT)

	tk.Plvsr(insertIDT, rowsAffectedT, errT)

	createStatementT := `
	CREATE TABLE TXTEST (
		ID NUMBER(10),
		USER_NAME VARCHAR(256),
		CODE NUMBER,
		UPDATE_TIME DATE
	)
	`

	insertIDT, rowsAffectedT, errT = ExecV(dbT, createStatementT)

	tk.CheckErrf("failed to create DB: %v", errT)

	tk.Plvsr(insertIDT, rowsAffectedT)

	insertIDT, rowsAffectedT, errT = ExecV(dbT, `insert into TXTEST (ID, USER_NAME, CODE, UPDATE_TIME) values(:a1, :a2, :a3, :a4)`, 1, "abc", 2.3, "2019-04-22T15:30:00")

	tk.CheckErrf("failed to insert record: %v", errT)

	tk.Plvsr(insertIDT, rowsAffectedT)

	tk.Plvsr(ExecV(dbT, `insert into TXTEST (ID, USER_NAME, CODE, UPDATE_TIME) values(:a1, :a2, :a3, :a4)`, 2, "abcd", 3.6, "2019-04-22T18:30:00"))

	tk.Plvsr(ExecV(dbT, `insert into TXTEST (USER_NAME, CODE, UPDATE_TIME) values(:a2, :a3, :a4)`, "abc", 1.2, "2019-04-25T19:30:00"))

	tk.Plvsr(ExecV(dbT, `insert into TXTEST (ID, CODE, UPDATE_TIME) values(:a1, :a3, :a4)`, 3, 1.2, "2019-04-22T15:30:00"))

	tk.Plvsr(ExecV(dbT, `insert into TXTEST (ID, USER_NAME, UPDATE_TIME) values(:a1, :a2, :a4)`, 4, "abc", time.Now().Format(time.RFC3339)))

	tk.Plvsr(ExecV(dbT, `insert into TXTEST (ID, USER_NAME, CODE) values(:a1, :a2, :a3)`, 5, "abc", 1.2))

	// QueryDBNS returns a SQL query result set, all values converted to string(empty string for NULL values)
	resultSetT, errT := QueryDBNS(dbT, "select * from TXTEST")

	tk.CheckErrf("failed to query: %v", errT)

	tk.Plvsr(resultSetT)

	// Parsing for DATE type values is a little bit complicated
	// below shows the examples for format such as 2019-04-22 15:30:00 +0000 UTC", "2019-04-23 14:24:44 +0800 +0800" and ""(empty or NULL value)
	timeStrT := resultSetT[4][3]

	timeT, errT := time.Parse("2006-01-02 15:04:05 -0700 MST", timeStrT)

	tk.Plvsr("converted time 1: ", timeT, tk.FormatTime(timeT, ""), errT)

	timeStrT = resultSetT[5][3]

	timeT, errT = time.Parse("2006-01-02 15:04:05 -0700 -0700", timeStrT)

	tk.Plvsr("converted time 2: ", timeT, tk.FormatTime(timeT, ""), errT)

	timeStrT = resultSetT[6][3]

	timeT, errT = time.Parse("", timeStrT)

	tk.Plvsr("converted time 3: ", timeT, tk.FormatTime(timeT, ""), errT)

	tk.Plvsr(QueryDBNS(dbT, "select count(*) from TXTEST"))

	resultT, errT := QueryDBI(dbT, "select * from TXTEST")

	tk.Plvsr(resultT, errT)

	timeT = (resultT[5][3]).(time.Time)

	tk.Pl("timeT: %v", timeT.UTC())
	tk.Pl("timeT: %v", timeT.Format(time.ANSIC))

	// QueryDBI is similar to QueryDBNS, but values are in interface{} type
	tk.Plvsr(QueryDBI(dbT, "select count(*) from TXTEST"))

	// passing parameters in SQL statement is supported
	tk.Plvsr(QueryDBNS(dbT, "select * from TXTEST where ID=?", 3))

}
