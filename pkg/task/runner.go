package task

import (
	"context"
	"fmt"
	"time"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/log"
)

type RunErrHandler func(item ScheduleItem, err error)

type SimpleRunner struct {
	*ScheduleTable
	RunErrorHandler RunErrHandler
}

func NewRunnerFromStore(store ScheduleTableStore) (*SimpleRunner, error) {
	table, err := store.GetFullTable()
	if err != nil {
		return nil, fmt.Errorf("get full table from store:%w", err)
	}
	return &SimpleRunner{
		ScheduleTable: &table,
	}, nil
}

func NewEmptyTaskRunner() *SimpleRunner {
	return &SimpleRunner{
		ScheduleTable: NewScheduleTable(),
	}
}

func NewSimpleRunnerWithItems(items ...*ScheduleItem) *SimpleRunner {
	return &SimpleRunner{
		ScheduleTable: NewScheduleTableWithItems(items...),
	}
}

func (r *SimpleRunner) Run(ctx context.Context) error {
	for _, item := range r.scheduleItemId2Desc {
		go r.loopOneItem(ctx, *item)
	}
	return nil
}

func (r *SimpleRunner) loopOneItem(ctx context.Context, item ScheduleItem) {
loop:
	for {
		now := time.Now()
		next, isOver := item.NextOccasion(now)
		if isOver {
			log.Infof("task [%s] over, quit loop")
			break loop
		}
		tm := time.NewTimer(next.Sub(now))
		select {
		case <-ctx.Done():
			log.Infof("SimpleRunner context done, quit loop")
			break loop
		case <-tm.C:
			//TODO: timeout
			r.executeOneItem(ctx, item)
		}
		tm.Stop()
	}
}

func (r *SimpleRunner) executeOneItem(ctx context.Context, item ScheduleItem) {
	if errExecute := item.Execute(ctx); errExecute != nil {
		log.Errorf("task [%s] execute failed:%s", errExecute.Error())
		if r.RunErrorHandler != nil {
			r.RunErrorHandler(item, errExecute)
		}
	}
}
