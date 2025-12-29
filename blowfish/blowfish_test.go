package blowfish

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBlowfish(t *testing.T) {
	Convey("blowfish", t, func() {
		bl, err := New(randomString(16), randomString(16))
		So(err, ShouldBeNil)
		data := []byte("hello world")
		enc, err := bl.Encrypt(data)
		So(err, ShouldBeNil)
		dec, err := bl.Decrypt(enc)
		So(err, ShouldBeNil)
		So(string(dec), ShouldEqual, string(data))
	})
}

const (
	uppercase  string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase         = "abcdefghijklmnopqrstuvwxyz"
	numeric           = "0123456789"
	alphabetic        = uppercase + lowercase + numeric
)

func randomString(length int) string {
	l := int64(len(alphabetic))
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = alphabetic[r.Int63()%l]
	}
	return string(b)
}
