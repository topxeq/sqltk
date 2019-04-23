# sqltk
Functions for simple SQL queries and commands.

Try to make SQL actions more easily.

Install:

go get -v github.com/topxeq/sqltk

Usage:

See the test file(sqltk_test.go) for details.

or in short:

	dbT, errT := sqltk.ConnectDB("goracle", dbConnectString)

	tk.CheckErrf("failed to connect: %v", errT)

	defer dbT.Close()

	insertIDT, rowsAffectedT, errT := sqltk.ExecV(dbT, `DROP TABLE TXTEST`)

	tk.Plvsr(insertIDT, rowsAffectedT, errT)

	createStatementT := `
	CREATE TABLE TXTEST (
		ID NUMBER(10),
		USER_NAME VARCHAR(256),
		CODE NUMBER,
		UPDATE_TIME DATE
	)
	`

	insertIDT, rowsAffectedT, errT = sqltk.ExecV(dbT, createStatementT)

	tk.CheckErrf("failed to create DB: %v", errT)

	tk.Plvsr(insertIDT, rowsAffectedT)

	insertIDT, rowsAffectedT, errT = sqltk.ExecV(dbT, `insert into TXTEST (ID, USER_NAME, CODE, UPDATE_TIME) values(:a1, :a2, :a3, TO_DATE(:a4 ,'yyyy-mm-dd hh24:mi:ss'))`, 1, "abc", 1.2, "2019-04-22 15:30:00")

	tk.CheckErrf("failed to insert record: %v", errT)

	tk.Plvsr(insertIDT, rowsAffectedT)

	tk.Plvsr(sqltk.ExecV(dbT, `insert into TXTEST (ID, USER_NAME, CODE, UPDATE_TIME) values(:a1, :a2, :a3, TO_DATE(:a4 ,'yyyy-mm-dd hh24:mi:ss'))`, 1, "abc", 1.2, "2019-04-22 15:30:00"))

	resultSetT, errT := sqltk.QueryDBNS(dbT, "select * from TXTEST")

	tk.CheckErrf("failed to query: %v", errT)

	tk.Plvsr(resultSetT)

	tk.Plvsr(sqltk.QueryDBNS(dbT, "select count(*) from TXTEST"))

	tk.Plvsr(sqltk.QueryDBI(dbT, "select * from TXTEST"))
