package pine

type ExpireError struct {
}

func NewExpireError() *ExpireError {
	return &ExpireError{}
}

func (e *ExpireError) Error() string {
	return "data expired"
}

type PageOutOfRangeError struct {
}

func NewPageOutOfRangeError() *PageOutOfRangeError {
	return &PageOutOfRangeError{}
}

func (e *PageOutOfRangeError) Error() string {
	return "page out of range"
}
