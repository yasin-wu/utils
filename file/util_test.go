package file

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCopy(t *testing.T) {
	Convey("blowfish", t, func() {
		err := Copy("./file.go", "./file_copy.go")
		defer os.RemoveAll("./file_copy.go")
		So(err, ShouldBeNil)
		err = os.MkdirAll("./internal_copy", 0755)
		So(err, ShouldBeNil)
		err = CopyDir("./internal", "./internal_copy")
		defer os.RemoveAll("./internal_copy")
		So(err, ShouldBeNil)
	})
}
