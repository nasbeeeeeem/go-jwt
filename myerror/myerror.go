package myerror

type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	return "Bad Request Error"
}

type InternalSeverError struct {
	Err error
}

func (e *InternalSeverError) Error() string {
	return "Internal Server Error"
}
