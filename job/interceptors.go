package job

import (
	"context"
	"fmt"
	"game-server/infra/log_bak"
)

type Work func(context.Context) error

func ErrorHandler(v interface{}) error {
	if err, ok := v.(error); ok {
		log_bak.Warnf("recover", "err_msg", err.Error())
		return err
	}
	returnErr := fmt.Errorf("Unknown Error, type: %T, value: %v", v, v)
	log_bak.Warnf("recover", "err_msg", returnErr.Error())
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
