package sqltk

import (
	"database/sql"
	"strings"

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
		dbT.Close()
		return nil, tk.Errf("failed to ping DB: %v", errT.Error())
	}

	return dbT, nil
}

// ConnectDBNoPing connected the database(with no ping action), don't forget to close it(probably by defer function)
func ConnectDBNoPing(driverStrA string, connectStrA string) (*sql.DB, error) {
	dbT, errT := sql.Open(driverStrA, connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
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

// QueryDBNSS execute a SQL query and return result set(first row will be the column names), all values will be string type(ensure for some DBs, such as MYSQL with uf8_general_ci encoding), can handle null values, passing parameters is supported as well.
func QueryDBNSS(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	resultSet := make([][]string, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	if columnSetT == nil {
		columnSetT, errT = rowsT.Columns()
		if errT != nil {
			return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
		}

		resultSet = append(resultSet, columnSetT)
	}

	for rowsT.Next() {
		rowCountT++

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

			resultRowS[k] = tk.Spr("%s", resultRow[k])
		}

		resultSet = append(resultSet, resultRowS)
	}

	errT = rowsT.Err()
	if errT != nil {
		return nil, tk.Errf("error occured while enumerating the result set: %v", errT.Error())
	}

	return resultSet, nil
}

// QueryDBNSV execute a SQL query and return result set(first row will be the column names), all values will be string type(ensure for some DBs, such as MYSQL with uf8_general_ci encoding), can handle null values, passing parameters is supported as well.
func QueryDBNSV(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
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

// QueryDBCount execute a SQL query for count(select count(*)), -1 indicates error, can handle null values, passing parameters is supported as well.
func QueryDBCount(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (int, error) {
	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return -1, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var countT int = -1

	for rowsT.Next() {
		errT = rowsT.Scan(&countT)
		if errT != nil {
			return -1, tk.Errf("failed to scan: %v", errT.Error())
		}

		break
	}

	return countT, nil
}

// QueryDBString execute a SQL query for a one string result, can handle null values, passing parameters is supported as well.
func QueryDBString(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (string, error) {
	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return "", tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var strT string = ""

	for rowsT.Next() {
		errT = rowsT.Scan(&strT)
		if errT != nil {
			return "", tk.Errf("failed to scan: %v", errT.Error())
		}

		return strT, nil
	}

	return "", tk.Errf("failed to get result: %v", "record not found")
}

// OneLineRecordToMap convert SQL result in [][]string (2 lines, first is the header) to map[string]string
func OneLineRecordToMap(recA [][]string) map[string]string {
	if recA == nil {
		return nil
	}

	if len(recA) < 2 {
		return nil
	}

	lenT := len(recA[0])

	mapT := make(map[string]string, lenT)

	for i := 0; i < lenT; i++ {
		mapT[recA[0][i]] = recA[1][i]
	}

	return mapT
}

// FormatSQLValue equivalent to strings.Replace(strA, "'", "''")
func FormatSQLValue(strA string) string {
	return strings.Replace(strA, "'", "''", -1)
}
