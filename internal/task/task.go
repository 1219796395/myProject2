package task

import (
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

// ProviderSet is task providers.
var ProviderSet = wire.NewSet(NewTask, NewNetworkConfigTask, NewRemoteConfigTask)

type Task struct {
	bc     *conf.Bootstrap
	log    *log.Helper
	cron   *cron.Cron
	ncTask *NetworkConfigTask

	rcTask *RemoteConfigTask
}

func NewTask(bc *conf.Bootstrap, logger log.Logger, ncTask *NetworkConfigTask, rcTask *RemoteConfigTask) *Task {
	t := &Task{
		bc:     bc,
		log:    log.NewHelper(logger),
		cron:   cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger))),
		ncTask: ncTask,
		rcTask: rcTask,
	}
	t.cron.AddFunc(
		t.bc.Biz.RemoteConfigCheckCacheByDbTask.Cron,
		func() { t.rcTask.CheckCacheByDB() },
	)
	t.cron.AddFunc(
		t.bc.Biz.NetworkConfigTranferStateTask.Cron,
		func() { t.ncTask.AutoTransferState() },
	)
	t.cron.Start()
	return t
}
