package mysql

import (
	"hcc/cello/lib/logger"
	"hcc/cello/lib/syscheck"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := syscheck.CheckRoot()
	if err != nil {
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer logger.FpLog.Close()

	err = Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer Db.Close()
}
