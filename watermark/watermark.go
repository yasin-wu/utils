package watermark

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/freetype"
)

type WaterMark struct {
}

type FontInfo struct {
	Size     float64 //文字大小
	Message  string  //文字内容
	Position int     //文字存放位置
	Dx       int     //文字x轴留白距离
	Dy       int     //文字y轴留白距离
	R        uint8   //文字颜色值RGBA中的R值
	G        uint8   //文字颜色值RGBA中的G值
	B        uint8   //文字颜色值RGBA中的B值
	A        uint8   //文字颜色值RGBA中的A值
}

func (w *WaterMark) New(srcFile, dstFile, fontFile string, fontInfo []FontInfo) error {
	imgFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	fileType := strings.Replace(path.Ext(path.Base(imgFile.Name())), ".", "", -1)
	switch fileType {
	case "gif":
		err = w.gifWaterMark(imgFile, dstFile, fontFile, fontInfo)
	default:
		err = w.staticWaterMark(imgFile, dstFile, fontFile, fileType, fontInfo)
	}
	return nil
}

func (this *WaterMark) staticWaterMark(srcFile *os.File, dstFile, fontFile, fileType string, fontInfo []FontInfo) error {
	var staticImg image.Image
	var err error
	switch fileType {
	case "png":
		staticImg, err = png.Decode(srcFile)
	default:
		staticImg, err = jpeg.Decode(srcFile)
	}
	if err != nil {
		return errors.New("image decode error: " + err.Error())
	}
	img := image.NewNRGBA(staticImg.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, staticImg.At(x, y))
		}
	}
	err = this.do(img, fontFile, fontInfo)
	if err != nil {
		return err
	}
	newFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer newFile.Close()
	switch fileType {
	case "png":
		err = png.Encode(newFile, img)
	default:
		err = jpeg.Encode(newFile, img, &jpeg.Options{Quality: 100})
	}
	return err
}

func (this *WaterMark) gifWaterMark(srcFile *os.File, dstFile, fontFile string, fontInfo []FontInfo) error {
	var err error
	gifImg, err := gif.DecodeAll(srcFile)
	if err != nil {
		return err
	}
	gifs := make([]*image.Paletted, 0)
	x0 := 0
	y0 := 0
	yuan := 0
	for k, v := range gifImg.Image {
		img := image.NewNRGBA(v.Bounds())
		if k == 0 {
			x0 = img.Bounds().Dx()
			y0 = img.Bounds().Dy()
		}
		if k == 0 && gifImg.Image[k+1].Bounds().Dx() > x0 && gifImg.Image[k+1].Bounds().Dy() > y0 {
			yuan = 1
			break
		}
		if x0 == img.Bounds().Dx() && y0 == img.Bounds().Dy() {
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					img.Set(x, y, v.At(x, y))
				}
			}
			err = this.do(img, fontFile, fontInfo)
			if err != nil {
				break
			}
			p1 := image.NewPaletted(v.Bounds(), v.Palette)
			draw.Draw(p1, v.Bounds(), img, image.ZP, draw.Src)
			gifs = append(gifs, p1)
		} else {
			gifs = append(gifs, v)
		}
	}
	if yuan == 1 {
		return errors.New("gif: image block is out of bounds")
	} else {
		if err != nil {
			return err
		}
		newFile, err := os.Create(dstFile)
		if err != nil {
			return err
		}
		defer newFile.Close()
		g1 := &gif.GIF{
			Image:     gifs,
			Delay:     gifImg.Delay,
			LoopCount: gifImg.LoopCount,
		}
		err = gif.EncodeAll(newFile, g1)
		return err
	}
}

func (this *WaterMark) do(img *image.NRGBA, fontFile string, fontInfo []FontInfo) error {
	var err error
	if fontFile == "" {
		fontFile = "./conf/captcha.ttf"
	}
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}
	errNum := 1
Loop:
	for _, v := range fontInfo {
		info := v.Message
		f := freetype.NewContext()
		f.SetDPI(108)
		f.SetFont(font)
		f.SetFontSize(v.Size)
		f.SetClip(img.Bounds())
		f.SetDst(img)
		f.SetSrc(image.NewUniform(color.RGBA{R: v.R, G: v.G, B: v.B, A: v.A}))
		first := 0
		two := 0
		switch v.Position {
		case TopLeft:
			first = v.Dx
			two = v.Dy + int(f.PointToFixed(v.Size)>>6)
		case TopRight:
			first = img.Bounds().Dx() - len(info)*4 - v.Dx
			two = v.Dy + int(f.PointToFixed(v.Size)>>6)
		case BottomLeft:
			first = v.Dx
			two = img.Bounds().Dy() - v.Dy
		case BottomRight:
			first = img.Bounds().Dx() - len(info)*4 - v.Dx
			two = img.Bounds().Dy() - v.Dy
		case Center:
			first = (img.Bounds().Dx() - len(info)*4) / 2
			two = (img.Bounds().Dy() - v.Dy) / 2
		default:
			errNum = 0
			break Loop
		}
		pt := freetype.Pt(first, two)
		_, err = f.DrawString(info, pt)
		if err != nil {
			break
		}
	}
	if errNum == 0 {
		err = errors.New("坐标值不对")
	}
	return err
}
