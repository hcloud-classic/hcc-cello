package end

import "hcc/cello/action/grpc/client"

// MainEnd : Main ending function
func MainEnd() {
	// rabbitmqEnd()
	mysqlEnd()
	loggerEnd()
	client.End()
}
