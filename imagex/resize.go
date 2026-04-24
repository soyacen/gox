package imagex

import (
	"golang.org/x/image/draw"
	"image"
)

// Resize returns a Transformer that resizes an image to the specified width and height.
//
// Parameters:
//   - width: the target width in pixels
//   - height: the target height in pixels
//   - scaler: the draw.Scaler used for scaling (e.g., draw.NearestNeighbor, draw.ApproxBiLinear, draw.CatmullRom)
//   - op: the composition operator (e.g., draw.Src, draw.Over)
//   - opts: optional scaling options; nil means default options
//
// Returns:
//   - Transformer: a function that performs the resize transformation
func Resize(width, height int, scaler draw.Scaler, op draw.Op, opts *draw.Options) Transformer {
	return func(img image.Image) image.Image {
		rect := image.Rect(0, 0, width, height)
		dst := image.NewRGBA(rect)
		scaler.Scale(dst, rect, img, img.Bounds(), op, opts)
		return dst
	}
}
