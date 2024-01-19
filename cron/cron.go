package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/robfig/cron/v3"
)

// Cron
// data写入内存
type Cron struct {
	data map[string][]cron.EntryID
	c    *cron.Cron
}

func New() *Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	c := cron.New(cron.WithParser(secondParser), cron.WithChain())
	c.Start()
	return &Cron{
		data: make(map[string][]cron.EntryID),
		c:    c,
	}
}

// Add Cron string 详见: https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.1#pkg-overview
// s m */h * * *  //每隔h小时m分s秒执行一次
// s m h * * ?    //每天h时m分钟s秒执行
// s m h * * w    //每星期wh时m分s秒执行,星期日-星期六(0-6)
// s m h d * ?    //每月d日h时m分s秒执行
// @every 1h      //每隔多久执行一次,1h30m10s
// id 需要唯一
func (c *Cron) Add(id, cron string, job cron.Job) (cron.EntryID, error) {
	if err := c.Delete(id); err != nil {
		return 0, err
	}
	entryID, err := c.c.AddJob(cron, job)
	if err != nil {
		return 0, err
	}
	c.updateData(id, entryID)
	return entryID, nil
}

// AddMany
// id 需要唯一,同一个任务的多次cron,全量增加
func (c *Cron) AddMany(id string, crons []string, job cron.Job) ([]cron.EntryID, error) {
	if err := c.Delete(id); err != nil {
		return nil, err
	}
	var entryIds []cron.EntryID
	for _, v := range crons {
		entryID, err := c.c.AddJob(v, job)
		if err != nil {
			return nil, err
		}
		entryIds = append(entryIds, entryID)
	}
	c.updateData(id, entryIds...)
	return entryIds, nil
}

func (c *Cron) Delete(id string) error {
	for _, v := range c.data[id] {
		c.c.Remove(v)
	}
	delete(c.data, id)
	return nil
}

func (c *Cron) DeleteAll() {
	for _, entry := range c.c.Entries() {
		c.c.Remove(entry.ID)
	}
	c.data = make(map[string][]cron.EntryID)
}

func (c *Cron) All() []cron.Entry {
	return c.c.Entries()
}

func (c *Cron) Data() map[string][]cron.EntryID {
	return c.data
}

func (c *Cron) updateData(id string, entryID ...cron.EntryID) {
	c.data[id] = append(c.data[id], entryID...)
}

// Every
// to @every
func Every(day, hour, minute, second int) string {
	if day+hour+minute+second == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString("@every ")
	buf.WriteString(strconv.Itoa(day*24 + hour))
	buf.WriteString("h")
	buf.WriteString(strconv.Itoa(minute))
	buf.WriteString("m")
	buf.WriteString(strconv.Itoa(second))
	buf.WriteString("s")
	return buf.String()
}

// ParseSchedule
// cron.Schedule to string
func ParseSchedule(schedule cron.Schedule) (string, error) {
	buf, err := json.Marshal(schedule)
	if err != nil {
		return "", err
	}
	switch schedule.(type) {
	case cron.ConstantDelaySchedule:
		var data cron.ConstantDelaySchedule
		err = json.Unmarshal(buf, &data)
		return data.Delay.String(), err
	case *cron.SpecSchedule:
		var data cron.SpecSchedule
		err = json.Unmarshal(buf, &data)
		return fmt.Sprintf("%d %d %d %d %d %d",
			data.Second, data.Minute, data.Hour, data.Dom, data.Month, data.Dow), err
	}
	return "", nil
}
