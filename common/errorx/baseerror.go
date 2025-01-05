package errorx

const defaultCode = 10000

type CodeError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code uint32, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// 为了迎合go当中的interface特性
func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) StatusCode() uint32 {
	return e.Code
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
