package errs

// Err provides a unified custom error used in the system, implements Error interface.
type Err struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Options struct {
	Message string
	Code    string
}

func New(opt Options) error {
	return &Err{
		Message: opt.Message,
		Code:    opt.Code,
	}
}

func (e *Err) Error() string {
	return e.Message
}

func IsCustom(e error) bool {
	_, ok := e.(*Err)
	return ok
}

func Code(err error) string {
	v, ok := err.(*Err)
	if !ok {
		return ""
	}

	return v.Code
}
