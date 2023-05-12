package job

import (
	"context"
	"testing"
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
