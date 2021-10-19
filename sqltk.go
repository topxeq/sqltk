package sqltk

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/topxeq/tk"
)

var versionG = "0.9a"

type SqlTK struct {
	Version string
}

var SqlTKX = &SqlTK{Version: versionG}

func (pA *SqlTK) NewSqlTK() *SqlTK {
	return &SqlTK{Version: versionG}
}

var NewSqlTK = SqlTKX.NewSqlTK

func (pA *SqlTK) GetVersion() string {
	return pA.Version
}

var GetVersion = SqlTKX.GetVersion

// ConnectDB connected the database, don't forget to close it(probably by defer function)
func (pA *SqlTK) ConnectDB(driverStrA string, connectStrA string) (*sql.DB, error) {
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

var ConnectDB = SqlTKX.ConnectDB

// ConnectDBNoPing connected the database(with no ping action), don't forget to close it(probably by defer function)
func (pA *SqlTK) ConnectDBNoPing(driverStrA string, connectStrA string) (*sql.DB, error) {
	dbT, errT := sql.Open(driverStrA, connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
	}

	return dbT, nil
}

var ConnectDBNoPing = SqlTKX.ConnectDBNoPing

// ExecV execute SQL statement, get the results(insert id and rows afftected), passing parameters is supported as well.
func (pA *SqlTK) ExecV(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (int64, int64, error) {
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

var ExecV = SqlTKX.ExecV

// QueryDBS execute a SQL query and return result set(first row will be the column names), all values will be string type, cannot handle null values, passing parameters is supported as well.
func (pA *SqlTK) QueryDBS(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
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

var QueryDBS = SqlTKX.QueryDBS

// QueryDBNS execute a SQL query and return result set(first row will be the column names), all values will be string type, can handle null values, passing parameters is supported as well.
func (pA *SqlTK) QueryDBNS(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
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

var QueryDBNS = SqlTKX.QueryDBNS

// QueryDBNSS execute a SQL query and return result set(first row will be the column names), all values will be string type(ensure for some DBs, such as MYSQL with uf8_general_ci encoding), can handle null values, passing parameters is supported as well.
func (pA *SqlTK) QueryDBNSS(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
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

var QueryDBNSS = SqlTKX.QueryDBNSS

// QueryDBNSSF the same as QueryDBNSS, but use special format on float values, format with argument floatFormatA(i.e. %1.2f etc).
func (pA *SqlTK) QueryDBNSSF(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
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

	colTypesT, errT := rowsT.ColumnTypes()
	if errT != nil {
		return nil, tk.Errf("failed to get column types of row %v: %v", rowCountT, errT.Error())
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
			// tk.Plvx(colTypesT[k].DatabaseTypeName())

			typeNameT := colTypesT[k].DatabaseTypeName()
			goTypeT := fmt.Sprintf("%T", resultRow[k])

			if tk.InStrings(typeNameT, "DOUBLE") {
				// resultRowS[k] = tk.Spr(floatFormatA, resultRow[k].(float64))
				resultRowS[k] = tk.Spr("%v", math.Round(tk.StrToFloat64(tk.Spr("%s", resultRow[k]), 0)*1000000)/1000000)
			} else if tk.InStrings(typeNameT, "NUMBER") && goTypeT == "int64" {
				tmps0 := tk.Spr("%v", resultRow[k])
				if tk.Contains(tmps0, ".") {
					tmps0 = strings.TrimRight(tmps0, "0")
				}

				if tk.EndsWith(tmps0, ".") {
					tmps0 = strings.TrimRight(tmps0, ".")
				}

				resultRowS[k] = tmps0
			} else if tk.InStrings(typeNameT, "DECIMAL", "NUMBER") {
				// tk.Pl("ROW: %v, %v", typeNameT, resultRow[k])
				tmps := tk.Spr("%s", resultRow[k])
				if tk.StartsWith(tmps, "%!s") {
					tk.Pl("ROW: %v, %T, %v", typeNameT, resultRow[k], resultRow[k])
					tmps = tk.Spr("%v", resultRow[k])
				}

				if tk.Contains(tmps, "e") {
					tmps = tk.Spr("%v", tk.ToInt(resultRow[k]))
				}

				if tk.Contains(tmps, ".") {
					tmps = strings.TrimRight(tmps, "0")
				}

				if tk.EndsWith(tmps, ".") {
					tmps = strings.TrimRight(tmps, ".")
				}

				resultRowS[k] = tmps
			} else if tk.InStrings(typeNameT, "INTEGER", "integer", "INT", "BIGINT") {
				tmps := tk.Spr("%v", resultRow[k])
				if tk.Contains(tmps, "[") {
					tmps = tk.ToStr(resultRow[k])
				}

				if tk.Contains(tmps, ".") {
					tmps = strings.TrimRight(tmps, "0")
				}

				if tk.EndsWith(tmps, ".") {
					tmps = strings.TrimRight(tmps, ".")
				}

				resultRowS[k] = tmps
			} else if tk.InStrings(typeNameT, "DATE") && goTypeT == "time.Time" {
				timeT, ok := resultRow[k].(time.Time)

				if ok {
					resultRowS[k] = tk.FormatTime(timeT)
				} else {
					resultRowS[k] = tk.Spr("%v", resultRow[k])
				}

			} else if tk.InStrings(typeNameT, "text", "TEXT", "CHAR", "VARCHAR", "VARCHAR2", "NVARCHAR2", "TIMESTAMP", "DATETIME") {
				resultRowS[k] = tk.Spr("%s", resultRow[k])
			} else {
				tk.Pl("ROW: %v, %T, %v", typeNameT, resultRow[k], resultRow[k])
				// tk.Pl("ROW: %v, %v", typeNameT, resultRow[k])
				resultRowS[k] = tk.Spr("%s", tk.ToStr(resultRow[k]))
			}
		}

		resultSet = append(resultSet, resultRowS)
	}

	errT = rowsT.Err()
	if errT != nil {
		return nil, tk.Errf("error occured while enumerating the result set: %v", errT.Error())
	}

	return resultSet, nil
}

var QueryDBNSSF = SqlTKX.QueryDBNSSF

// QueryDBNSV execute a SQL query and return result set(first row will be the column names), all values will be string type(ensure for some DBs, such as MYSQL with uf8_general_ci encoding), can handle null values, passing parameters is supported as well.
func (pA *SqlTK) QueryDBNSV(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]string, error) {
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

var QueryDBNSV = SqlTKX.QueryDBNSV

// QueryDBI execute a SQL query and return result set(first row will be the column names), all values will be interface{} type, passing parameters is supported as well.
func (pA *SqlTK) QueryDBI(dbA *sql.DB, sqlStrA string, argsA ...interface{}) ([][]interface{}, error) {

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

var QueryDBI = SqlTKX.QueryDBI

// QueryDBCount execute a SQL query for count(select count(*)), -1 indicates error, can handle null values, passing parameters is supported as well. Also used to get a single int result from SQL query.
func (pA *SqlTK) QueryDBCount(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (int, error) {
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

var QueryDBCount = SqlTKX.QueryDBCount

// QueryDBFloat execute a SQL query for get a single float value, can handle null values, passing parameters is supported as well.
func (pA *SqlTK) QueryDBFloat(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (float64, error) {
	rowsT, errT := dbA.Query(sqlStrA, argsA...)

	if errT != nil {
		return 0, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var countT float64 = 0

	for rowsT.Next() {
		errT = rowsT.Scan(&countT)
		if errT != nil {
			return 0, tk.Errf("failed to scan: %v", errT.Error())
		}

		break
	}

	countT = math.Round(countT*1000000) / 1000000

	return countT, nil
}

var QueryDBFloat = SqlTKX.QueryDBFloat

// QueryDBString execute a SQL query for a one string result, can handle null values, passing parameters is supported as well.
func (pA *SqlTK) QueryDBString(dbA *sql.DB, sqlStrA string, argsA ...interface{}) (string, error) {
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

var QueryDBString = SqlTKX.QueryDBString

// OneLineRecordToMap convert SQL result in [][]string (2 lines, first is the header) to map[string]string
func (pA *SqlTK) OneLineRecordToMap(recA [][]string) map[string]string {
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

var OneLineRecordToMap = SqlTKX.OneLineRecordToMap

// OneColumnRecordsToArray convert SQL result in [][]string (several lines, first is the header, only one column per line) to []string
func (pA *SqlTK) OneColumnRecordsToArray(recsA [][]string) []string {
	if recsA == nil {
		return nil
	}

	if len(recsA) < 1 {
		return nil
	}

	lenT := len(recsA)

	aryT := make([]string, 0, lenT)

	for i := 0; i < lenT; i++ {
		if i == 0 {
			continue
		}
		aryT = append(aryT, recsA[i][0])
	}

	return aryT
}

var OneColumnRecordsToArray = SqlTKX.OneColumnRecordsToArray

func (pA *SqlTK) RecordsToMapArrayMap(recA [][]string, keyA string) map[string][]map[string]string {
	if recA == nil {
		return nil
	}

	lenT := len(recA)

	if lenT < 1 {
		return nil
	}

	return tk.TableToMSSMapArray(recA, keyA)
}

var RecordsToMapArrayMap = SqlTKX.RecordsToMapArrayMap

// RecordsToMapArray convert SQL result in [][]string (multi lines, first is the header) to []map[string]string
func (pA *SqlTK) RecordsToMapArray(recA [][]string) []map[string]string {
	if recA == nil {
		return nil
	}

	lenT := len(recA)

	if lenT < 1 {
		return nil
	}

	lineLenT := len(recA[0])

	aryT := make([]map[string]string, lenT-1)

	for i := 1; i < lenT; i++ {
		mapT := make(map[string]string, lenT)

		for j := 0; j < lineLenT; j++ {
			mapT[recA[0][j]] = recA[i][j]
		}

		aryT[i-1] = mapT
	}

	return aryT
}

var RecordsToMapArray = SqlTKX.RecordsToMapArray

// FormatSQLValue equivalent to strings.Replace(strA, "'", "''")
func (pA *SqlTK) FormatSQLValue(strA string) string {
	strT := strings.Replace(strA, "\r", "\\r", -1)
	strT = strings.Replace(strT, "\n", "\\n", -1)
	strT = strings.Replace(strT, "'", "''", -1)

	return strT
}

var FormatSQLValue = SqlTKX.FormatSQLValue

func (pA *SqlTK) ConnectDBX(driverStrA string, connectStrA string) interface{} {
	dbT, errT := ConnectDBNoPing(driverStrA, connectStrA)

	if errT != nil {
		return errT
	}

	return dbT
}

var ConnectDBX = SqlTKX.ConnectDBX

func (pA *SqlTK) ExecDBX(dbA *sql.DB, sqlStrA string, argsA ...interface{}) interface{} {
	idT, affectT, errT := ExecV(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	return []int64{idT, affectT}
}

var ExecDBX = SqlTKX.ExecDBX

func (pA *SqlTK) QueryDBX(dbA *sql.DB, sqlStrA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBNSSF(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	if len(sqlRsT) < 1 {
		return tk.Errf("invalid record length")
	}

	return tk.TableToMSSArray(sqlRsT)
}

var QueryDBX = SqlTKX.QueryDBX

func (pA *SqlTK) QueryDBRecsX(dbA *sql.DB, sqlStrA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBNSSF(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	if len(sqlRsT) < 1 {
		return tk.Errf("invalid record length")
	}

	return sqlRsT
}

var QueryDBRecsX = SqlTKX.QueryDBRecsX

func (pA *SqlTK) QueryDBMapX(dbA *sql.DB, sqlStrA string, idA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBNSSF(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	if len(sqlRsT) < 1 {
		return tk.Errf("invalid record length")
	}

	return tk.TableToMSSMap(sqlRsT, idA)
}

var QueryDBMapX = SqlTKX.QueryDBMapX

func (pA *SqlTK) QueryDBMapArrayX(dbA *sql.DB, sqlStrA string, idA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBNSSF(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	if len(sqlRsT) < 1 {
		return tk.Errf("invalid record length")
	}

	return tk.TableToMSSMapArray(sqlRsT, idA)
}

var QueryDBMapArrayX = SqlTKX.QueryDBMapArrayX

func (pA *SqlTK) QueryCountX(dbA *sql.DB, sqlStrA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBCount(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	if sqlRsT < 0 {
		return tk.Errf("result error: %v", sqlRsT)
	}

	return sqlRsT
}

var QueryCountX = SqlTKX.QueryCountX

func (pA *SqlTK) QueryFloatX(dbA *sql.DB, sqlStrA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBFloat(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	return sqlRsT
}

var QueryFloatX = SqlTKX.QueryFloatX

func (pA *SqlTK) QueryStringX(dbA *sql.DB, sqlStrA string, argsA ...interface{}) interface{} {
	sqlRsT, errT := QueryDBCount(dbA, sqlStrA, argsA...)

	if errT != nil {
		return errT
	}

	return sqlRsT
}

var QueryStringX = SqlTKX.QueryStringX

func (pA *SqlTK) CloseDBX(dbA *sql.DB) error {

	return dbA.Close()
}

var CloseDBX = SqlTKX.CloseDBX
