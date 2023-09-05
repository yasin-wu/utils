package strings

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelector(t *testing.T) {
	Convey("selector match", t, func() {
		ok, err := Match("(tag1 || tag2) && (tag3 && tag4) || tag5", []string{"tag2", "tag3", "tag4"})
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)
		ok, err = Match("(tag1 || tag2) && (tag3 && tag4) || (tag5)", []string{"tag2", "tag3", "tag5"})
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)
		ok, err = Match("(tag1 || tag2) && (tag3 && tag4) || (tag5)", []string{"tag2", "tag3", "tag6"})
		So(err, ShouldBeNil)
		So(ok, ShouldBeFalse)
	})
}
