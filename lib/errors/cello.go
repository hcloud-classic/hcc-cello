package errors

const (
	CelloInternalInitFail                 = cello + internal + initFail
	CelloInternalOperationFail            = cello + internal + operationFail
	CelloInternalUUIDGenerationError      = cello + internal + UUIDGenerationError
	CelloInternalTimeStampConversionError = cello + internal + timestampConversionError
	CelloInternalCreateServerFailed       = cello + internal + createServerFailed
	CelloInternalCreateServerRoutineError = cello + internal + createServerRoutineError
	CelloInternalGetAvailableNodesError   = cello + internal + getAvailableNodesError
	CelloInternalServerNodePresentError   = cello + internal + serverNodePresentError

	CelloInternalVolumeHandleError = cello + internal + volumehandlerError

	CelloInternalStoragePoolError  = cello + internal + storagepoolError
	CelloInternalCreateVolumeError = cello + internal + createvolumeError
	CelloInternalPreparePxeError   = cello + internal + preparepxeError
	CelloInternalWriteIscsiError   = cello + internal + writeiscsiError

	CelloGrpcArgumentError = cello + grpc + argumentError
	CelloGrpcRequestError  = cello + grpc + requestError
	CelloGrpcGetNodesError = cello + grpc + getNodesError

	CelloSQLOperationFail = cello + sql + operationFail
	CelloSQLNoResult      = cello + sql + noResult
	CelloSQLArgumentError = cello + sql + argumentError
)
