package myerror

type BadRequestError struct {
	Err error
}

// http code 400
func (e *BadRequestError) Error() string {
	return "Bad Request Error"
}

// http code 500
type InternalSeverError struct {
	Err error
}

func (e *InternalSeverError) Error() string {
	return "Internal Server Error"
}
