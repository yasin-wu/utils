package mfa

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/skip2/go-qrcode"

	"github.com/yasin-wu/utils/mfa/internal/aes"
)

const (
	uppercase      string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase             = "abcdefghijklmnopqrstuvwxyz"
	alphabetic            = uppercase + lowercase
	numeric               = "0123456789"
	alphanumeric          = alphabetic + numeric
	cacheKeyPrefix        = "mfa:"
)

type MFA struct {
	product   string
	secretKey string
	cli       *redis.Client
	redisExp  time.Duration
	digits    int
}

type QrCode struct {
	Key       string
	Code      string
	PNG       []byte
	PNGBase64 string
}

type Recovery struct {
	Key      string   `json:"key"`
	Recovery []string `json:"recovery"`
	StandBy  []string `json:"stand_by"`
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func New(product, secretKey string, conf *RedisConfig) (*MFA, error) {
	if conf == nil {
		return nil, errors.New("redis config is nil")
	}
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &MFA{
		product:   product,
		secretKey: secretKey,
		cli:       cli,
		redisExp:  10 * time.Minute,
		digits:    6,
	}, nil
}

func (m *MFA) QrCode(username string) (*QrCode, error) {
	random := m.randomString(10)
	if err := m.cli.Set(context.Background(), cacheKeyPrefix+username, random, m.redisExp).Err(); err != nil {
		return nil, err
	}
	dst := make([]byte, 16)
	base32.StdEncoding.Encode(dst, []byte(random))
	key := string(dst)
	qr := &QrCode{
		Key:  key,
		Code: "otpauth://totp/" + m.product + "-" + username + "?secret=" + key,
	}
	png, err := qrcode.Encode(qr.Code, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	qr.PNG = png
	qr.PNGBase64 = base64.StdEncoding.EncodeToString(png)
	return qr, nil
}

func (m *MFA) Recovery(username, code string) (*Recovery, error) {
	cacheKey, err := m.cli.Get(context.Background(), cacheKeyPrefix+username).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("cacheKey:", cacheKey)
	now := time.Now()
	if code != m.totp([]byte(cacheKey), now) &&
		code != m.totp([]byte(cacheKey), now.Add(30*time.Second)) &&
		code != m.totp([]byte(cacheKey), now.Add(-30*time.Second)) {
		return nil, errors.New("mfa authentication failed")
	}
	skey := m.skey(username)
	standBy, recovery := m.standBy(skey)
	dst := make([]byte, 16)
	base32.StdEncoding.Encode(dst, []byte(cacheKey))
	rev := &Recovery{
		Key:      string(aes.Encode(dst, []byte(skey))),
		Recovery: recovery,
		StandBy:  standBy,
	}
	return rev, nil
}

func (m *MFA) Check(username, key, code string, recovery ...string) ([]string, bool, error) {
	var (
		newRecovery []string
		passed      bool
	)
	skey := m.skey(username)
	for _, v := range recovery {
		if code == string(aes.Decode([]byte(v), []byte(skey))) {
			passed = true
		} else {
			newRecovery = append(newRecovery, v)
		}
	}
	encodeKey, err := base32.StdEncoding.DecodeString(string(aes.Decode([]byte(key), []byte(skey))))
	if err != nil {
		return nil, false, err
	}
	now := time.Now()
	if code == m.totp(encodeKey, now) ||
		code == m.totp(encodeKey, now.Add(30*time.Second)) ||
		code == m.totp(encodeKey, now.Add(-30*time.Second)) {
		passed = true
	}
	return newRecovery, passed, nil
}

func (m *MFA) randomString(length int) string {
	l := int64(len(alphanumeric))
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = alphanumeric[r.Int63()%l]
	}
	return string(b)
}

func (m *MFA) totp(key []byte, t time.Time) string {
	counter := uint64(t.UnixNano()) / 30e9
	h := hmac.New(sha1.New, key)
	_ = binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < m.digits && i < 8; i++ {
		d *= 10
	}
	return fmt.Sprintf("%0*d", m.digits, int(v%d))
}

func (m *MFA) standBy(skey string) ([]string, []string) {
	var (
		standby  []string
		recovery []string
	)
	for i := 0; i < 10; i++ {
		src := m.randomString(10)
		standby = append(standby, src)
		s := string(aes.Encode([]byte(src), []byte(skey)))
		recovery = append(recovery, s)
	}
	return standby, recovery
}

func (m *MFA) skey(username string) string {
	return string(aes.Decode([]byte(m.secretKey), []byte(m.secretKey))) + username
}
