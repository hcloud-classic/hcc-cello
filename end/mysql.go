package end

import "hcc/cello/lib/mysql"

func mysqlEnd() {
	if mysql.Db != nil {
		_ = mysql.Db.Close()
	}
}
