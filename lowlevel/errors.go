package lowlevel

import "errors"

var (
	ErrIncorrectHeader              = errors.New("incorrect header")
	ErrIncorrectTag                 = errors.New("incorrect tag")
	ErrIncorrectValueRepresentation = errors.New("incorrect value representation")
	ErrIncorrectValueLength         = errors.New("incorrect value length")
	ErrIncorrectUInt32              = errors.New("incorrect uint32")
)
