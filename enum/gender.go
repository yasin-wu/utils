//go:generate stringer -type Gender -linecomment
package enum

// Gender go generate ./
// Gender 性别
type Gender int

const (
	GenderMale   Gender = 1 //男
	GenderFemale Gender = 2 //女
)

func (c Gender) Int() int {
	return int(c)
}

func (c Gender) Equals(i int) bool {
	return c.Int() == i
}
