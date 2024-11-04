package zip

import (
	"testing"
)

func TestZip(t *testing.T) {
	err := Compresses("/Users/yasin/Downloads/20241104", "./tmp/20241104.zip")
	if err != nil {
		t.Fatal(err)
	}
	err = Compresses("/Users/yasin/Downloads/20241104", "./tmp/20241104-password.zip", "123456")
	if err != nil {
		t.Fatal(err)
	}
}
