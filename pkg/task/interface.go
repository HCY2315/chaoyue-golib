//Package task 实现基于时间的任务执行, 类似cron
package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/microservice"

	"github.com/google/uuid"
)

type Task interface {
	Name() string
	Execute(ctx context.Context) error
}

type ScheduleItemDesc interface {
	microservice.PersistObject
	// NextOccasion返回下次时间和是否结束
	NextOccasion(timePoint time.Time) (time.Time, bool)
}

type intervalScheduleItem struct {
	CreateTime time.Time
	Interval   time.Duration
}

func NewIntervalScheduleItem(interval time.Duration) *intervalScheduleItem {
	if interval < 0 {
		panic(fmt.Sprintf("need positive interval"))
	}
	return &intervalScheduleItem{Interval: interval, CreateTime: time.Now()}
}

func (i intervalScheduleItem) Encode() []byte {
	bytes, _ := json.Marshal(i)
	return bytes
}

func (i *intervalScheduleItem) Recover(persistBytes []byte) error {
	err := json.Unmarshal(persistBytes, i)
	if err != nil {
		return fmt.Errorf("recover intervalScheduleItem from bytes %s:%w", persistBytes, err)
	}
	return nil
}

func (i intervalScheduleItem) NextOccasion(timePoint time.Time) (time.Time, bool) {
	if timePoint.Before(i.CreateTime) {
		return i.CreateTime, false
	}
	cycleNum := timePoint.Sub(i.CreateTime)/i.Interval + 1
	next := i.CreateTime.Add(i.Interval * cycleNum)
	return next, false
}

//ScheduleItemDesc 表示一条任务的调度信息
type ScheduleItem struct {
	ScheduleItemDesc
	Task
	TaskId string
}

func (si ScheduleItem) String() string {
	return fmt.Sprintf("ScheduleItem:[%s(%s)]", si.Name(), si.TaskId)
}

type ScheduleTable struct {
	//id -> desc
	scheduleItemId2Desc map[string]*ScheduleItem
}

func (st *ScheduleTable) AddItem(desc ScheduleItemDesc, task Task) error {
	id := uuid.New().String()
	st.scheduleItemId2Desc[id] = &ScheduleItem{
		ScheduleItemDesc: desc,
		Task:             task,
		TaskId:           id,
	}
	return nil
}

func (st ScheduleTable) ForEach(invoker func(taskId string, item ScheduleItem)) {
	for taskId, item := range st.scheduleItemId2Desc {
		invoker(taskId, *item)
	}
}

func NewScheduleTable() *ScheduleTable {
	return &ScheduleTable{
		scheduleItemId2Desc: make(map[string]*ScheduleItem),
	}
}

func NewScheduleTableWithItems(items ...*ScheduleItem) *ScheduleTable {
	descMap := make(map[string]*ScheduleItem, len(items))
	for _, item := range items {
		descMap[item.TaskId] = item
	}
	return &ScheduleTable{scheduleItemId2Desc: descMap}
}

//ScheduleTableStore 存储/读取/恢复 任务表
type ScheduleTableStore interface {
	SaveOneItem(ScheduleItem) error
	GetFullTable() (ScheduleTable, error)
}
