package job

import (
	"context"
	"fmt"

	"github.com/cr-mao/lori/log"
)

type Work func(context.Context) error

func ErrorHandler(v interface{}) error {
	if err, ok := v.(error); ok {
		log.Errorw("recover_err", err.Error())
		return err
	}
	returnErr := fmt.Errorf("Unknown Error, type: %T, value: %v", v, v)
	log.Errorw("recover_err", returnErr.Error())
	return returnErr
}

func RecoverInterceptor(next Work) Work {
	return func(ctx context.Context) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = ErrorHandler(r)
			}
		}()
		return next(ctx)
	}
}
