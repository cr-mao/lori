package job

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestOnceJob(t *testing.T) {
	var a int
	onceJob := Once{
		Work: func(ctx context.Context) error {
			a = 1
			return nil
		},
	}
	onceJob.Run(context.Background())
	if a != 1 {
		t.Log("not ok")
	} else {
		t.Log("ok")
	}
}

/*
once job
ERROR recover_err=Unknown Error, type: string, value: throw  exception
    once_test.go:38: end
*/
func TestInterceptorsUnkownError(t *testing.T) {
	onceJob := Once{
		Work: func(ctx context.Context) error {
			fmt.Println("once job")
			panic("throw  exception")
			return nil
		},
	}
	onceJob.Run(context.Background())

	time.Sleep(1 * time.Second)
	t.Log("end")
}

/*
once job
ERROR recover_err=test error
    once_test.go:51: end
*/
func TestInterceptorsWithError(t *testing.T) {
	onceJob := Once{
		Work: func(ctx context.Context) error {
			fmt.Println("once job")
			panic(errors.New("test error"))
			return nil
		},
	}
	onceJob.Run(context.Background())
	time.Sleep(1 * time.Second)
	t.Log("end")
}
