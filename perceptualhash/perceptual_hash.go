// Package perceptualhash provides utilities for computing a perceptual hash of images.
package perceptualhash

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"

	"golang.org/x/image/draw"
)

// Config holds debugging options for perceptual hashing.
type Config struct {
	Debug          bool
	DebugParameter struct {
		PreprocessedImagePath string
		VisualizedImagePath   string
	}
}

var defaultConfig = Config{
	Debug: false,
}

var ErrUnsupportedFormat = errors.New("image format is not supported")

// FromPath computes the perceptual hash of the image at filePath.
// It optionally accepts a custom configuration.
func FromPath(filePath string, configs ...Config) (string, error) {
	// load the configurations
	config := defaultConfig
	if len(configs) > 0 {
		config = configs[0]
	}

	// 1. Load the image
	loadedImage, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer loadedImage.Close()

	// 2. Decode the image
	decodedImage, format, err := image.Decode(loadedImage)
	if err != nil {
		return "", err
	}

	if format != "png" && format != "jpeg" && format != "jpg" {
		return "", ErrUnsupportedFormat
	}

	// 3. Preprocess the image
	preprocessedImage := preprocessImage(decodedImage, config)
	if config.Debug {
		if err := saveImage(preprocessedImage, format, config.DebugParameter.PreprocessedImagePath); err != nil {
			return "", err
		}
	}

	hash := generateHash(preprocessedImage)
	if config.Debug {
		if err := visualizeHash(hash, format, config); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%016x", hash), nil
}

// CompareHashes compares two perceptual hashes and returns the Hamming distance.
// The distance is the number of differing bits between the two hashes.
func CompareHashes(hash1, hash2 string) (int, error) {
	if len(hash1) != len(hash2) {
		return 0, fmt.Errorf("hashes must be of the same length")
	}

	distance := 0
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			distance++
		}
	}

	return distance, nil
}

// preprocessImage resizes the image to 32x32 and converts it to grayscale.
func preprocessImage(inputImage image.Image, config Config) *image.Gray {
	resizedImage := image.NewGray(image.Rect(0, 0, 32, 32))
	draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), inputImage, inputImage.Bounds(), draw.Over, nil)

	return resizedImage
}

// saveImage writes the given image to location in the specified format.
func saveImage(img image.Image, format string, location string) error {
	outputImage, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer outputImage.Close()

	switch format {
	case "png":
		err = png.Encode(outputImage, img)
		if err != nil {
			return err
		}
	case "jpg", "jpeg":
		err = jpeg.Encode(outputImage, img, &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// generateHash computes the DCT-based 64-bit hash from a 32x32 grayscale image.
func generateHash(img *image.Gray) uint64 {
	var pixels [][]float64
	for y := 0; y < 32; y++ {
		row := make([]float64, 32)
		for x := 0; x < 32; x++ {
			grayColor := img.GrayAt(x, y)
			row[x] = float64(grayColor.Y)
		}
		pixels = append(pixels, row)
	}

	dctMatrix := dct(pixels)
	var dctValues []float64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			dctValues = append(dctValues, dctMatrix[y][x])
		}
	}

	var sum float64
	for i := 1; i < len(dctValues); i++ {
		sum += dctValues[i]
	}
	average := sum / 63

	var hash uint64
	for i, value := range dctValues {
		if i > 0 && value > average {
			hash |= 1 << i
		}
	}

	return hash
}

// dct performs a 2D Discrete Cosine Transform on the input matrix.
func dct(matrix [][]float64) [][]float64 {
	N := len(matrix)
	dct := make([][]float64, N)
	for u := 0; u < N; u++ {
		dct[u] = make([]float64, N)
		for v := 0; v < N; v++ {
			sum := 0.0
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					sum += matrix[x][y] *
						math.Cos((float64(2*x+1)*float64(u)*math.Pi)/(2*float64(N))) *
						math.Cos((float64(2*y+1)*float64(v)*math.Pi)/(2*float64(N)))
				}
			}

			cu := 1.0
			cv := 1.0
			if u == 0 {
				cu = 1 / math.Sqrt2
			}
			if v == 0 {
				cv = 1 / math.Sqrt2
			}
			dct[u][v] = 0.25 * cu * cv * sum
		}
	}
	return dct
}

// visualizeHash creates a small 8x8 image from hash bits for debugging.
func visualizeHash(hash uint64, format string, config Config) error {
	size := 8
	img := image.NewGray(image.Rect(0, 0, size, size))
	for i := range size {
		for j := range size {
			bitPosition := uint(i*size + j)
			bit := (hash >> bitPosition) & 1
			var pixelColor color.Gray

			if bit == 1 {
				pixelColor = color.Gray{255}
			} else {
				pixelColor = color.Gray{0}
			}

			img.SetGray(j, i, pixelColor)
		}
	}
	return saveImage(img, format, config.DebugParameter.VisualizedImagePath)
}
