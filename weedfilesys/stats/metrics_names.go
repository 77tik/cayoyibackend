package stats

const (
	// master client
	FailedToKeepConnected = "failedToKeepConnected"
	FailedToSend          = "failedToSend"
	FailedToReceive       = "failedToReceive"
	RedirectedToLeader    = "redirectedToLeader"
	OnPeerUpdate          = "onPeerUpdate"
	Failed                = "failed"

	// volume server
	WriteToLocalDisk   = "writeToLocalDisk"
	WriteToReplicas    = "writeToReplicas"
	DownloadLimitCond  = "downloadLimitCondition"
	UploadLimitCond    = "uploadLimitCondition"
	ReadProxyReq       = "readProxyRequest"
	ReadRedirectReq    = "readRedirectRequest"
	EmptyReadProxyLoc  = "emptyReadProxyLocaction"
	FailedReadProxyReq = "failedReadProxyRequest"

	ErrorSizeMismatchOffsetSize = "errorSizeMismatchOffsetSize"
	ErrorSizeMismatch           = "errorSizeMismatch"
	ErrorCRC                    = "errorCRC"
	ErrorIndexOutOfRange        = "errorIndexOutOfRange"
	ErrorGetNotFound            = "errorGetNotFound"
	ErrorGetInternal            = "errorGetInternal"
)
