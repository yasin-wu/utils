package mfa

import (
	"testing"

	"github.com/skip2/go-qrcode"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	username  = "yasin"
	product   = "Aizen"
	secretKey = "7ix51vt0whithoq"
	redisConf = &RedisConfig{
		Addr:     "127.0.0.1:6379",
		Username: "",
		Password: "",
		DB:       10,
	}
)

func TestQrCode(t *testing.T) {
	Convey("qrcode", t, func() {
		m, err := New(product, secretKey, redisConf)
		So(err, ShouldBeNil)
		qr, err := m.QrCode(username)
		So(err, ShouldBeNil)
		err = qrcode.WriteFile(qr.Code, qrcode.Medium, 256, "./qr.png")
		So(err, ShouldBeNil)
	})
}

func TestRecovery(t *testing.T) {
	Convey("recovery", t, func() {
		m, err := New(product, secretKey, redisConf)
		So(err, ShouldBeNil)
		rev, err := m.Recovery(username, "309071")
		So(err, ShouldBeNil)
		t.Log(rev.Recovery)
		t.Log(rev.StandBy)
		t.Log(rev.Key)
	})
}

func TestCheck(t *testing.T) {
	Convey("check", t, func() {
		m, err := New(product, secretKey, redisConf)
		So(err, ShouldBeNil)
		recovery := []string{
			"83GUmhPfvJ0kKxCJeSL0j46dbheuaqoYFpY=",
			"r/KgKcA1wPSte3Zw13+xGdhbNWMtlUqE90o=",
			"KfpEavXLNnN2II4yE+sN0YxuPwm5jXOkbec=",
		}
		key := "hVx8SZuqfllxY3s/nrd+KGkCh6ug3k5tPmQxY7+BRo8="
		//standBy := []string{
		//	"vPjGBMNYkU",
		//	"5pPb1OvGt6",
		//	"CS2ANh11iI",
		//}
		rev, ok, err := m.Check(username, key, "CS2ANh11iI", recovery...)
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)
		t.Log(rev)
	})
}
