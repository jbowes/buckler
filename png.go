package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"code.google.com/p/freetype-go/freetype"
)

func makePngShield(w http.ResponseWriter, d Data) {
	w.Header().Add("content-type", "image/png")

	fi, _ := os.Open("edge.png")
	edge, _ := png.Decode(fi)
	mask := image.NewRGBA(image.Rect(0, 0, 100, 19))
	draw.Draw(mask, edge.Bounds(), edge, image.ZP, draw.Src)

	sr := image.Rect(2, 0, 3, 19)
	for i := 3; i <= 97; i++ {
		dp := image.Point{i, 0}
		r := sr.Sub(sr.Min).Add(dp)
		draw.Draw(mask, r, edge, sr.Min, draw.Src)
	}

	sr = image.Rect(0, 0, 1, 19)
	dp := image.Point{99, 0}
	r := sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, edge, sr.Min, draw.Src)

	sr = image.Rect(1, 0, 2, 19)
	dp = image.Point{98, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, edge, sr.Min, draw.Src)

	img := image.NewRGBA(image.Rect(0, 0, 100, 19))
	right := color.RGBA{69, 203, 20, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{right}, image.ZP, draw.Src)

	left := color.RGBA{79, 79, 79, 255}
	rect := image.Rect(0, 0, 50, 19)
	draw.Draw(img, rect, &image.Uniform{left}, image.ZP, draw.Src)

	dst := image.NewRGBA(image.Rect(0, 0, 100, 19))
	draw.DrawMask(dst, dst.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

	fontBytes, err := ioutil.ReadFile("opensanssemibold.ttf")
	if err != nil {
		log.Println(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(10)
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)

	pt := freetype.Pt(6, 13)
	offset, _ := c.DrawString("build", pt)

	pt = freetype.Pt(53, 13)
	c.DrawString("passing", pt)

	println(offset.X, offset.Y)
	png.Encode(w, dst)
}
