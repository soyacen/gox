package imagex

import "image"

// Rotate90 rotates the given image 90 degrees clockwise.
//
// Parameters:
//   - oldImage: the source image to rotate
//
// Returns:
//   - image.Image: a new RGBA image rotated 90 degrees clockwise
func Rotate90(oldImage image.Image) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, oldImage.Bounds().Dy(), oldImage.Bounds().Dx()))
	for x := oldImage.Bounds().Min.Y; x < oldImage.Bounds().Max.Y; x++ {
		for y := oldImage.Bounds().Max.X - 1; y >= oldImage.Bounds().Min.X; y-- {
			newImage.Set(oldImage.Bounds().Max.Y-x, y, oldImage.At(y, x))
		}
	}
	return newImage
}

// Rotate180 rotates the given image 180 degrees.
//
// Parameters:
//   - oldImage: the source image to rotate
//
// Returns:
//   - image.Image: a new RGBA image rotated 180 degrees
func Rotate180(oldImage image.Image) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, oldImage.Bounds().Dx(), oldImage.Bounds().Dy()))
	for x := oldImage.Bounds().Min.X; x < oldImage.Bounds().Max.X; x++ {
		for y := oldImage.Bounds().Min.Y; y < oldImage.Bounds().Max.Y; y++ {
			newImage.Set(oldImage.Bounds().Max.X-x, oldImage.Bounds().Max.Y-y, oldImage.At(x, y))
		}
	}
	return newImage
}

// Rotate270 rotates the given image 270 degrees clockwise (or 90 degrees counter-clockwise).
//
// Parameters:
//   - oldImage: the source image to rotate
//
// Returns:
//   - image.Image: a new RGBA image rotated 270 degrees clockwise
func Rotate270(oldImage image.Image) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, oldImage.Bounds().Dy(), oldImage.Bounds().Dx()))
	for x := oldImage.Bounds().Min.Y; x < oldImage.Bounds().Max.Y; x++ {
		for y := oldImage.Bounds().Max.X - 1; y >= oldImage.Bounds().Min.X; y-- {
			newImage.Set(x, oldImage.Bounds().Max.X-y, oldImage.At(y, x))
		}
	}
	return newImage
}
