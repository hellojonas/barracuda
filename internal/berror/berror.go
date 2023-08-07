package berror

const (
	ErrSourceNotReachable    = 1
	ErrSourceStatusNotOK     = 2
	ErrSourceParseFailed     = 3
	ErrDBInvalidQuery        = 4
	ErrDBQueryFailed         = 5
	ErrDBTxInintFailed       = 6
	ErrEmptyArticle          = 7
	ErrDBResultExtractFailed = 8
	ErrSourceNotFound        = 9
)

type BError struct {
	code    int
	message string
}

func New(code int, msg string) BError {
	return BError{}
}

func (e BError) Error() string {
	return e.message
}
