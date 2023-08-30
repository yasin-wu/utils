package ldap

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLdap(t *testing.T) {
	Convey("ldap", t, func() {
		l := New("127.0.0.1:389", "administrator", "yasin123", "dc=y,dc=w,dc=u,dc=com")
		result, err := l.SearchGroup()
		So(err, ShouldBeNil)
		fmtPrint(result[0])
		fmtPrint(result[1])
		presult, err := l.SearchPerson()
		So(err, ShouldBeNil)
		fmtPrint(presult[0])
		fmtPrint(presult[1])
	})
}

func fmtPrint(data any) {
	buffer, _ := json.MarshalIndent(data, "", "")
	fmt.Println(string(buffer))
}
