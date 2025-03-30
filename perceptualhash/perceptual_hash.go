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

// preprocessImage convert image to grayscale and convert it to 32x32
func preprocessImage(inputImage image.Image, config Config) *image.Gray {
	resizedImage := image.NewGray(image.Rect(0, 0, 32, 32))
	draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), inputImage, inputImage.Bounds(), draw.Over, nil)

	return resizedImage
}

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
