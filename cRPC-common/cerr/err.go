package cerr

type CXError struct {
	ErrCode int64
	ErrMsg  string
}

func (ce *CXError) Error() string {
	return ce.ErrMsg
}

func NewError(errCode int64, errMsg string) *CXError {
	return &CXError{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}

var (
	ErrDB              = NewError(3001, "database error")
	ErrInvalidParam    = NewError(3002, "invalid param")
	ErrTimeOut         = NewError(3003, "timeout")
	ErrAuth            = NewError(3004, "permission denied")
	ErrBusy            = NewError(3005, "server busy")
	ErrInternal        = NewError(3006, "internal error")
	ErrRedisNil        = NewError(3007, "redis: nil")
	ErrHttpHeaderErr   = NewError(3008, "header err")
	ErrServiceNotFound = NewError(3404, "service not find")
	ErrConnNotFound    = NewError(3405, "conn not find")
)
