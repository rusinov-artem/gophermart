package error

type Error struct {
	Error appError
}

type appError interface {
	appError()
}
