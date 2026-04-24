package imagex

import "image"

// Transformer is a function type that transforms an image.
//
// Parameters:
//   - image.Image: the source image to transform
//
// Returns:
//   - image.Image: the transformed image
type Transformer func(image.Image) image.Image

// Transform applies a sequence of Transformer functions to an image.
// Transformers are applied in reverse order (last one first).
//
// Parameters:
//   - img: the source image to transform
//   - transformers: the transformations to apply, applied from last to first
//
// Returns:
//   - image.Image: the transformed image
func Transform(img image.Image, transformers ...Transformer) image.Image {
	for i := len(transformers) - 1; i >= 0; i-- {
		img = transformers[i](img)
	}
	return img
}
