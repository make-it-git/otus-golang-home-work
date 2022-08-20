package logic

import "errors"

type ErrBussinessLogic struct {
	Err error
}

func (e ErrBussinessLogic) Error() string {
	return e.Err.Error()
}

var (
	ErrMissingID              = NewBusinessLogicError(errors.New("id not provided"))
	ErrDuplicateID            = NewBusinessLogicError(errors.New("id duplicate"))
	ErrNotFoundID             = NewBusinessLogicError(errors.New("id not found"))
	ErrEndTimeBeforeStartTime = NewBusinessLogicError(errors.New("end time should be after start time"))
)

func NewBusinessLogicError(err error) ErrBussinessLogic {
	return ErrBussinessLogic{Err: err}
}
