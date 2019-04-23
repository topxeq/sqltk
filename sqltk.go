package sqltk

import (
	"database/sql"

	"github.com/topxeq/tk"
)

// ConnectDB connected the database, don't forget to close it(probably by defer function)
func ConnectDB(driverStrA string, connectStrA string) (*sql.DB, error) {
	dbT, errT := sql.Open(driverStrA, connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
	}

	errT = dbT.Ping()

	if errT != nil {
		return nil, tk.Errf("failed to ping DB: %v", errT.Error())
	}

	return dbT, nil
}

// ExecV execute SQL statement, get the results(insert id and rows afftected), passing parameters is supported as well.
func ExecV(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (int64, int64, error) {
	resultT, errT := dbA.Exec(sqlStrA, argsA...)
	if errT != nil {
		return 0, 0, tk.Errf("failed to exec: %v", errT.Error())
	}

	insertIDT, errT := resultT.LastInsertId()

	if errT != nil {
		insertIDT = 0
		// return 0, 0, tk.Errf("failed to get result insertID: %v", errT.Error())
	}

	rowAffectedT, errT := resultT.RowsAffected()

	if errT != nil {
		rowAffectedT = 0
		// return 0, 0, tk.Errf("failed to get result rowAffected: %v", errT.Error())
	}

	return insertIDT, rowAffectedT, nil

}

// QueryDBS execute a SQL query and return result set(first row will be the column names), all values will be string type, cannot handle null values, passing parameters is supported as well.
func QueryDBS(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	resultSet := make([][]string, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	for rowsT.Next() {
		rowCountT++

		if columnSetT == nil {
			columnSetT, errT = rowsT.Columns()
			if errT != nil {
				return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
			}

			resultSet = append(resultSet, columnSetT)
		}

		columnLenT := len(columnSetT)
		var resultRow = make([]string, columnLenT)
		var resultRowP = make([]interface{}, columnLenT)

		for k := 0; k < columnLenT; k++ {
			resultRowP[k] = &(resultRow[k])
		}

		errT = rowsT.Scan(resultRowP...)
		if errT != nil {
			return nil, tk.Errf("failed to scan %v: %v", rowCountT, errT.Error())
		}

		resultSet = append(resultSet, resultRow)
	}

	errT = rowsT.Err()
	if errT != nil {
		return nil, tk.Errf("error occured while enumerating the result set: %v", errT.Error())
	}

	return resultSet, nil
}

// QueryDBNS execute a SQL query and return result set(first row will be the column names), all values will be string type, can handle null values, passing parameters is supported as well.
func QueryDBNS(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	resultSet := make([][]string, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	for rowsT.Next() {
		rowCountT++

		if columnSetT == nil {
			columnSetT, errT = rowsT.Columns()
			if errT != nil {
				return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
			}

			resultSet = append(resultSet, columnSetT)
		}

		columnLenT := len(columnSetT)

		var resultRow = make([]interface{}, columnLenT)
		var resultRowP = make([]interface{}, columnLenT)
		var resultRowS = make([]string, columnLenT)

		for k := 0; k < columnLenT; k++ {
			resultRowP[k] = &(resultRow[k])
		}

		errT = rowsT.Scan(resultRowP...)
		if errT != nil {
			return nil, tk.Errf("failed to scan %v: %v", rowCountT, errT.Error())
		}

		for k := 0; k < columnLenT; k++ {
			if resultRow[k] == nil {
				resultRowS[k] = ""
				continue
			}

			resultRowS[k] = tk.Spr("%v", resultRow[k])
		}

		resultSet = append(resultSet, resultRowS)
	}

	errT = rowsT.Err()
	if errT != nil {
		return nil, tk.Errf("error occured while enumerating the result set: %v", errT.Error())
	}

	return resultSet, nil
}

// QueryDBI execute a SQL query and return result set(first row will be the column names), all values will be interface{} type, passing parameters is supported as well.
func QueryDBI(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]interface{}, error) {

	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	resultSet := make([][]interface{}, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	for rowsT.Next() {
		rowCountT++

		if columnSetT == nil {
			columnSetT, errT = rowsT.Columns()
			if errT != nil {
				return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
			}

			lenT := len(columnSetT)
			setT := make([]interface{}, lenT)
			for k := 0; k < lenT; k++ {
				setT[k] = columnSetT[k]
			}

			resultSet = append(resultSet, setT)

		}

		columnLenT := len(columnSetT)

		var resultRow = make([]interface{}, columnLenT)
		var resultRowP = make([]interface{}, columnLenT)

		for k := 0; k < columnLenT; k++ {
			resultRowP[k] = &(resultRow[k])
		}

		errT = rowsT.Scan(resultRowP...)
		if errT != nil {
			return nil, tk.Errf("failed to scan %v: %v", rowCountT, errT.Error())
		}

		resultSet = append(resultSet, resultRow)
	}

	errT = rowsT.Err()
	if errT != nil {
		return nil, tk.Errf("error occured while enumerating the result set: %v", errT.Error())
	}

	return resultSet, nil
}