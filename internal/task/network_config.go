package task

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// 自定义的定时任务
type NetworkConfigTask struct {
	bc    *conf.Bootstrap
	log   *log.Helper
	logic *biz.NetworkConfigLogic
}

// NetworkConfigTask new a timer task.
func NewNetworkConfigTask(bc *conf.Bootstrap, logger log.Logger, logic *biz.NetworkConfigLogic) *NetworkConfigTask {
	return &NetworkConfigTask{
		bc:    bc,
		log:   log.NewHelper(logger),
		logic: logic,
	}
}

// 开始执行定时任务
func (t *NetworkConfigTask) AutoTransferState() error {
	if !t.bc.Biz.NetworkConfigTranferStateTask.Switch {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.bc.Biz.NetworkConfigTranferStateTask.LockExpire)*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.log.WithContext(ctx).Errorf("network_config cron task timeout")
		return nil
	default:
	}
	if err := t.logic.AutoTransferState(ctx); err != nil {
		t.log.WithContext(ctx).Errorf("transfer network_config state fail! err = %+v", err)
	}
	return nil
}
