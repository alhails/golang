package exception

type UnsupportedOperationException struct {
	*RuntimeException
}

func NewUnsupportedOperationException() *UnsupportedOperationException {
	return &UnsupportedOperationException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewUnsupportedOperationException1(message string) *UnsupportedOperationException {
	return &UnsupportedOperationException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewUnsupportedOperationException2(message string, cause Throwable) *UnsupportedOperationException {
	return &UnsupportedOperationException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
