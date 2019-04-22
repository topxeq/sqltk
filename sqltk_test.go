package sqltk

import (
	"testing"

	"github.com/topxeq/tk"
	_ "gopkg.in/goracle.v2"
)

var dbConnectString = `test/test@127.0.0.1:1521/testdb`

func Test001(t *testing.T) {
	rs, error := SelectDBS(dbConnectString, "select * from TEST")

	tk.Plvsr(rs, error)

}
