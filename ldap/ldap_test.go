package ldap

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yasin-wu/utils/util"
	"testing"
)

func TestLdap(t *testing.T) {
	Convey("ldap", t, func() {
		l := New("127.0.0.1:389", "Administrator", "yasinwu", "DC=winServer2008R2,DC=com")
		entries, err := l.Search(PersonFilter, 100)
		So(err, ShouldBeNil)
		t.Log(len(entries))
		util.PrintlnFmt(entries)
		err = l.Add(PersonClass, "yasinwu")
		So(err, ShouldBeNil)
	})
}
