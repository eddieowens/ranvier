package except

type ErrorBundle interface {
	error
	Errors() []error
	AddError(e error)
}
