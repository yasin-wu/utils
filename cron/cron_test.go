package cron

import (
	"fmt"
	"testing"

	"github.com/robfig/cron/v3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAddCronJob(t *testing.T) {
	Convey("add cron job", t, func() {
		cronJob := New()
		cronTest := &CornTest{ID: "cron_test"}
		entryID, err := cronJob.Add(cronTest.ID, "@daily", cronTest)
		So(err, ShouldBeNil)
		So(cronJob.data[cronTest.ID], ShouldResemble, []cron.EntryID{entryID})
		entries := cronJob.All()
		So(len(entries), ShouldEqual, 1)
		err = cronJob.Delete(cronTest.ID)
		So(err, ShouldBeNil)
		entries = cronJob.All()
		So(len(entries), ShouldEqual, 0)
	})
}

func TestAddManyCronJob(t *testing.T) {
	Convey("add many cron job", t, func() {
		cronJob := New()
		cronTest := &CornTest{ID: "cron_test"}
		crons := []string{"@daily", "@yearly", "@monthly"}
		entryIds, err := cronJob.AddMany(cronTest.ID, crons, cronTest)
		So(err, ShouldBeNil)
		So(cronJob.data[cronTest.ID], ShouldResemble, entryIds)
		entries := cronJob.All()
		So(len(entries), ShouldEqual, 3)
		crons = []string{"@yearly", "@monthly"}
		entryIds, err = cronJob.AddMany(cronTest.ID, crons, cronTest)
		So(err, ShouldBeNil)
		So(cronJob.data[cronTest.ID], ShouldResemble, entryIds)
		entries = cronJob.All()
		So(len(entries), ShouldEqual, 2)
		err = cronJob.Delete(cronTest.ID)
		So(err, ShouldBeNil)
		entries = cronJob.All()
		So(len(entries), ShouldEqual, 0)
	})
}

type CornTest struct {
	ID string `json:"id"`
}

func (c *CornTest) Run() {
	fmt.Println("this is test cron")
}

func TestCronEvery(t *testing.T) {
	Convey("cron every", t, func() {
		So(Every(0, 1, 2, 3), ShouldEqual, "@every 1h2m3s")
		So(Every(0, 0, 2, 3), ShouldEqual, "@every 0h2m3s")
		So(Every(1, 1, 2, 3), ShouldEqual, "@every 25h2m3s")
		So(Every(0, 0, 0, 0), ShouldEqual, "")
	})
}
