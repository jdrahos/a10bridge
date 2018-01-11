package api

type A10Error interface {
	Code() int
	Message() string
	Error() string
}
