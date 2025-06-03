package ldap

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLdap(t *testing.T) {
	convey.Convey("ldap", t, func() {
		l := New("127.0.0.1:389", "administrator", "yasin123", "dc=y,dc=w,dc=u,dc=com")
		units, err := l.SearchUnit()
		convey.So(err, convey.ShouldBeNil)
		fmt.Println("unit len: ", len(units))
		fmtPrint(units[0])
		fmtPrint(units[1])
		fmtPrint(units[2])
		persons, err := l.SearchPerson()
		convey.So(err, convey.ShouldBeNil)
		fmt.Println("person len: ", len(persons))
		fmtPrint(persons[0])
		fmtPrint(persons[1])
	})
}

func fmtPrint(data any) {
	buffer, _ := json.MarshalIndent(data, "", "")
	fmt.Println(string(buffer))
}
