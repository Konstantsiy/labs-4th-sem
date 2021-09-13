package gamma

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func clamp(f float64) uint8 {
	v := int64(f+0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}

func fillColor(c color.RGBA, filler []uint8) color.RGBA {
	return color.RGBA{
		R: filler[c.R],
		G: filler[c.G],
		B: filler[c.B],
		A: filler[c.A],
	}
}

func asRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	res := image.NewRGBA(bounds)
	draw.Draw(res, bounds, src, bounds.Min, draw.Src)
	return res
}

func AddGamma(img image.Image, c, y float64) *image.RGBA {
	filler := make([]uint8, 256)

	for i := 0; i < 256; i++ {
		filler[i] = clamp(c + math.Pow(float64(i) / 255, y) * 255.0)
	}

	bounds := img.Bounds()
	res := asRGBA(img)
	w, h := bounds.Dx(), bounds.Dy()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			resPos := y*res.Stride+x*4

			c := color.RGBA{}

			dr := &res.Pix[resPos+0]
			dg := &res.Pix[resPos+1]
			db := &res.Pix[resPos+2]
			da := &res.Pix[resPos+3]

			c.R = *dr
			c.G = *dg
			c.B = *db
			c.A = *da

			c = fillColor(c, filler)

			*dr = c.R
			*dg = c.G
			*db = c.B
			*da = c.A
		}
	}

	return res
}
