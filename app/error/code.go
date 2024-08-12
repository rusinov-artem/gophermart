package error

type Code int

var (
	LoginAlreadyInUse        Code = 1
	UnableToCheckLoginExists Code = 2
	UnableToSaveUser         Code = 3
	UnableToSaveToken        Code = 4
	InvalidCredentials       Code = 5
	ServiceUnavailable       Code = 6
	OrderNrExists            Code = 7
	BadOrderOwnership        Code = 8
)
