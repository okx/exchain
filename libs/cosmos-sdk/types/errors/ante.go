package errors

const (
	AnteErrorsLowerBoundary = 100000
)

// AnteError if the error occurs in ante. It's a AnteError
type anteError struct {
	Err error
}

func NewAnteError(err error) *anteError {
	return &anteError{
		Err: err,
	}
}

func (ae *anteError) Error() string {
	return ae.Err.Error()
}

func (ae *anteError) ABCICode() uint32 {
	return abciCode(ae.Err) + AnteErrorsLowerBoundary
}
