package sqltk

import (
	"database/sql"

	"github.com/topxeq/tk"
	_ "gopkg.in/goracle.v2"
)

// RecordCell hold each cell in result set for SQL query
// Type: 0 - unknown, 1 - string, 2 - int, 3 - float64, 4 - datetime
type RecordCell struct {
	Name   string
	Type   string
	IsNull bool
	Value  string
}

func (p *RecordCell) GetIntValue(defaultA int) int {
	return tk.StrToIntWithDefaultValue(p.Value, defaultA)
}

func SelectDB(connectStrA string, sqlStrA string) ([][]RecordCell, error) {
	dbT, errT := sql.Open("goracle", connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
	}

	defer dbT.Close()

	rowsT, errT := dbT.Query(sqlStrA)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var columnSet []RecordCell = nil
	var columnCountA = 0

	// var resultSet = make([][]RecordCell, 0)

	for rowsT.Next() {

		if columnSet == nil {
			columnInfoT, errT := rowsT.ColumnTypes()

			if errT != nil {
				return nil, tk.Errf("failed to retrieve column info: %v", errT.Error())
			}

			columnSet = make([]RecordCell, 0)

			for _, v := range columnInfoT {
				isNullableT, okT := v.Nullable()

				if !okT {
					isNullableT = false
				}

				precisionT, scaleT, okT := v.DecimalSize()

				var precisionStr = ""

				if okT {
					precisionStr = tk.Spr("%v.%v", precisionT, scaleT)
				}

				rcT := RecordCell{Name: v.Name(), Type: v.DatabaseTypeName(), IsNull: isNullableT, Value: precisionStr}

				columnSet = append(columnSet, rcT)
				// tk.Pl("%v: %#v", i, v)

				// tk.Pl("scanType: %#v", v.ScanType().Kind())
			}

			tk.Pl("columnSet: %#v", columnSet)
		}

		columnCountA = len(columnSet)

		for j := 0; j < columnCountA; j++ {

		}

		// rowsT.ColumnTypes

		// errT = rowsT.Scan(&id, &userID)
		// if errT != nil {
		// 	return nil, tk.Errf("遍历查询结果时发生错误：%v", errT.Error())
		// }

		// sb.WriteString(tk.Spr("id: %v, userID: %v\n", id, userID))
	}

	errT = rowsT.Err()
	if errT != nil {
		return nil, tk.Errf("查询结果有错误：%v", errT.Error())
	}

	return nil, nil
}

func SelectDBS(connectStrA string, sqlStrA string) ([][]string, error) {
	dbT, errT := sql.Open("goracle", connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
	}

	defer dbT.Close()

	rowsT, errT := dbT.Query(sqlStrA)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var resultSet [][]string = make([][]string, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	// var valueT string

	for rowsT.Next() {
		rowCountT++

		if columnSetT == nil {
			columnSetT, errT = rowsT.Columns()
			if errT != nil {
				return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
			}
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

func SelectDBI(connectStrA string, sqlStrA string) ([][]interface{}, error) {
	dbT, errT := sql.Open("goracle", connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
	}

	defer dbT.Close()

	rowsT, errT := dbT.Query(sqlStrA)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var resultSet [][]interface{} = make([][]interface{}, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	// var valueT string

	for rowsT.Next() {
		rowCountT++

		if columnSetT == nil {
			columnSetT, errT = rowsT.Columns()
			if errT != nil {
				return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
			}
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

func SelectDBVI(connectStrA string, sqlStrA string, argsA ...interface{}) ([][]interface{}, error) {
	dbT, errT := sql.Open("goracle", connectStrA)

	if errT != nil {
		return nil, tk.Errf("failed to open DB: %v", errT.Error())
	}

	defer dbT.Close()

	stmtT, errT := dbT.Prepare(sqlStrA)
	if errT != nil {
		return nil, tk.Errf("failed to prepare SQL statement: %v", errT.Error())
	}

	defer stmtT.Close()

	rowsT, errT := stmtT.Query(argsA...)

	if errT != nil {
		return nil, tk.Errf("failed to run query: %v", errT.Error())
	}

	defer rowsT.Close()

	var resultSet [][]interface{} = make([][]interface{}, 0)
	var rowCountT = 0
	var columnSetT []string = nil

	// var valueT string

	for rowsT.Next() {
		rowCountT++

		if columnSetT == nil {
			columnSetT, errT = rowsT.Columns()
			if errT != nil {
				return nil, tk.Errf("failed to get columns of row %v: %v", rowCountT, errT.Error())
			}
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
