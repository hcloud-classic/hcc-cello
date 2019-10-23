package mysql

import (
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"
	"hcc/cello/lib/syscheck"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	if !syscheck.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = Db.Close()
	}()
}
