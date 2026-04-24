package imagex

import (
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
)

// Png decodes a PNG image from src, applies the given transformers, and encodes the result to dest as PNG.
//
// Parameters:
//   - dest: the writer to encode the transformed PNG image to
//   - src: the reader to decode the PNG image from
//   - transformers: optional image transformations to apply
//
// Returns:
//   - error: any error encountered during decode, transform, or encode
func Png(dest io.Writer, src io.Reader, transformers ...Transformer) error {
	img, err := png.Decode(src)
	if err != nil {
		return err
	}
	return png.Encode(dest, Transform(img, transformers...))
}

// Jpeg decodes a JPEG image from src, applies the given transformers, and encodes the result to dest as JPEG.
//
// Parameters:
//   - dest: the writer to encode the transformed JPEG image to
//   - src: the reader to decode the JPEG image from
//   - opts: JPEG encoding options; nil means default options
//   - transformers: optional image transformations to apply
//
// Returns:
//   - error: any error encountered during decode, transform, or encode
func Jpeg(dest io.Writer, src io.Reader, opts *jpeg.Options, transformers ...Transformer) error {
	img, err := jpeg.Decode(src)
	if err != nil {
		return err
	}
	return jpeg.Encode(dest, Transform(img, transformers...), opts)
}

// Gif decodes a GIF image from src, applies the given transformers, and encodes the result to dest as GIF.
//
// Parameters:
//   - dest: the writer to encode the transformed GIF image to
//   - src: the reader to decode the GIF image from
//   - opt: GIF encoding options; nil means default options
//   - transformers: optional image transformations to apply
//
// Returns:
//   - error: any error encountered during decode, transform, or encode
func Gif(dest io.Writer, src io.Reader, opt *gif.Options, transformers ...Transformer) error {
	img, err := gif.Decode(src)
	if err != nil {
		return err
	}
	return gif.Encode(dest, Transform(img, transformers...), opt)
}

// Bmp decodes a BMP image from src, applies the given transformers, and encodes the result to dest as BMP.
//
// Parameters:
//   - dest: the writer to encode the transformed BMP image to
//   - src: the reader to decode the BMP image from
//   - transformers: optional image transformations to apply
//
// Returns:
//   - error: any error encountered during decode, transform, or encode
func Bmp(dest io.Writer, src io.Reader, transformers ...Transformer) error {
	img, err := bmp.Decode(src)
	if err != nil {
		return err
	}
	return bmp.Encode(dest, Transform(img, transformers...))
}

// Webp decodes a WebP image from src, applies the given transformers, and encodes the result to dest as PNG.
//
// Parameters:
//   - dest: the writer to encode the transformed PNG image to
//   - src: the multipart file to decode the WebP image from
//   - transformers: optional image transformations to apply
//
// Returns:
//   - error: any error encountered during decode, transform, or encode
func Webp(dest io.Writer, src multipart.File, transformers ...Transformer) error {
	img, err := webp.Decode(src)
	if err != nil {
		return err
	}
	return png.Encode(dest, Transform(img, transformers...))
}
