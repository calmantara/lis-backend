package errors

import (
	"github.com/rotisserie/eris"
)

func New(msg string) error {
	return eris.New(msg)
}

func Wrapf(err error, msg string, args ...any) error {
	return eris.Wrapf(err, msg, args...)
}

func Wrap(err error, msg string) error {
	return eris.Wrap(err, msg)
}

func Is(err, target error) bool {
	return eris.Is(err, target)
}

func As(err error, target any) bool {
	return eris.As(err, target)
}

func Unpack(err error) (res []string) {
	unpacked := eris.Unpack(err)
	for _, e := range unpacked.ErrChain {
		res = append(res, e.Msg)
	}
	res = append(res, unpacked.ErrRoot.Msg)

	return res
}
