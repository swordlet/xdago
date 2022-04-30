package consensus

import (
	"fmt"
	"xdago/core"
)

type Task struct {
	task      core.XdagField
	taskTime  uint64
	taskIndex uint64
}

func (t *Task) Task() core.XdagField {
	return t.task
}

func (t *Task) SetTask(task core.XdagField) {
	t.task = task
}

func (t *Task) TaskTime() uint64 {
	return t.taskTime
}

func (t *Task) SetTaskTime(taskTime uint64) {
	t.taskTime = taskTime
}

func (t *Task) TaskIndex() uint64 {
	return t.taskIndex
}

func (t *Task) SetTaskIndex(taskIndex uint64) {
	t.taskIndex = taskIndex
}

func (t *Task) ToString() string {
	return fmt.Sprintf("Task: { taskTime: %d, %016x }", t.taskTime, t.taskTime)
}
