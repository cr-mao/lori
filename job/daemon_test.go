package job

import (
	"context"
	"testing"
	"time"
)

func TestDaemonJob(t *testing.T) {
	var result = 0
	go func() {
		daemonJob := Daemon{
			Work: func(ctx context.Context) error {
				result++
				return nil
			},
			Rate: time.Second * 5,
		}
		daemonJob.Run(context.Background())
	}()

	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*31)
	defer cancel()
END:
	for {
		select {
		case <-ctx.Done():
			t.Logf("result end is %d", result)
			t.Logf("err +%s", ctx.Err())
			break END
		case <-ticker.C:
			t.Logf("result now is %d", result)
		}
	}
	//最终 result 基本都是6 。 每5秒这样会加1.
}
