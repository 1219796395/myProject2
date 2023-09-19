package task

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// 自定义的定时任务
type RemoteConfigTask struct {
	bc    *conf.Bootstrap
	log   *log.Helper
	logic *biz.RemoteConfigLogic
}

// NetworkConfigTask new a timer task.
func NewRemoteConfigTask(bc *conf.Bootstrap, logger log.Logger, logic *biz.RemoteConfigLogic) *RemoteConfigTask {
	// 初始化remote_config的本地cache
	log.Infof("start init remote_config cache")
	if err := logic.CheckCacheByDB(context.Background()); err != nil {
		panic(err)
	}
	log.Infof("success init remote_config cache !")
	return &RemoteConfigTask{
		bc:    bc,
		log:   log.NewHelper(logger),
		logic: logic,
	}
}

func (t *RemoteConfigTask) CheckCacheByDB() error {
	if !t.bc.Biz.RemoteConfigCheckCacheByDbTask.Switch {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.bc.Biz.RemoteConfigCheckCacheByDbTask.LockExpire)*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.log.WithContext(ctx).Errorf("network_config cron task timeout")
		return nil
	default:
	}
	if err := t.logic.CheckCacheByDB(ctx); err != nil {
		t.log.WithContext(ctx).Errorf("transfer network_config state fail! err = %+v", err)
	}
	return nil
}
