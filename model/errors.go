package app

// set NonRetryableErrorTypes
type BadRequestError struct{}

func (m *BadRequestError) Error() string {
	return "Request is invalid"
}
