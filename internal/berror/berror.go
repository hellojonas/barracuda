package berror

const (
	ErrSourceNotReachable = 1
	ErrSourceStatusNotOK  = 2
	ErrSourceParseError   = 3
)

type bError struct {
	code    int
	message string
}

func New(code int, msg string) bError {
	return bError{}
}

func (e bError) Error() string {
	return e.message
}
