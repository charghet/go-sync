package web

type RecordExistsErr struct {
	Msg string
}

func (e *RecordExistsErr) Error() string {
	return e.Msg
}

type ServiceErr struct {
	Code int
	Msg  string
}

func (e ServiceErr) Error() string {
	return e.Msg
}

type InnerErr struct {
	Code int
	Msg  string
	Err  error
}

func (e InnerErr) Error() string {
	return e.Msg
}
