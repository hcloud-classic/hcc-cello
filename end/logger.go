package end

import "hcc/cello/lib/logger"

func loggerEnd() {
	_ = logger.FpLog.Close()
}
