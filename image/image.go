package image

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
	"unsafe"

	"github.com/yasin-wu/utils/consts"
)

const _png = "png"

func FileToBase64(file string) (string, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()
	fileType := strings.ReplaceAll(path.Ext(path.Base(imgFile.Name())), ".", "")
	var staticImg image.Image
	switch fileType {
	case _png:
		staticImg, err = png.Decode(imgFile)
	default:
		staticImg, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		return "", err
	}
	emptyBuff := bytes.NewBuffer(nil)
	switch fileType {
	case _png:
		_ = png.Encode(emptyBuff, staticImg)
	default:
		_ = jpeg.Encode(emptyBuff, staticImg, nil)
	}
	dist := make([]byte, 20*consts.MB)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	index := bytes.IndexByte(dist, 0)
	baseImage := dist[0:index]
	return *(*string)(unsafe.Pointer(&baseImage)), nil
}

func ImageToBase64(img image.Image, fileType string) (string, error) {
	var err error
	emptyBuff := bytes.NewBuffer(nil)
	switch fileType {
	case _png:
		err = png.Encode(emptyBuff, img)
	default:
		err = jpeg.Encode(emptyBuff, img, nil)
	}
	if err != nil {
		return "", err
	}
	dist := make([]byte, 20*1024*1024)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	index := bytes.IndexByte(dist, 0)
	baseImage := dist[0:index]
	return *(*string)(unsafe.Pointer(&baseImage)), nil
}
