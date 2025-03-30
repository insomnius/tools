package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/insomnius/tools/perceptualhash"
)

func main() {
	// Path to the folder containing images
	imagesFolder := "./images"

	// Output file to save the results
	outputFile, err := os.Create("hashes.txt")
	if err != nil {
		log.Fatal("Failed to create output file:", err)
	}
	defer outputFile.Close()

	// Walk through the images folder
	err = filepath.Walk(imagesFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-image files
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		// Create debug folder if needed
		debugFolder := filepath.Join(filepath.Dir(path), "debug")
		if _, err := os.Stat(debugFolder); os.IsNotExist(err) {
			os.Mkdir(debugFolder, 0755)
		}

		// Configure perceptual hash with debug options
		conf := perceptualhash.Config{
			Debug: true,
		}

		if conf.Debug {
			baseName := filepath.Base(path)
			conf.DebugParameter.PreprocessedImagePath = filepath.Join(debugFolder, "preprocessed_"+baseName)
			conf.DebugParameter.VisualizedImagePath = filepath.Join(debugFolder, "visualized_"+baseName)
		}

		// Generate hash
		hash, err := perceptualhash.FromPath(path, conf)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", path, err)
			return nil // Continue with next file
		}

		// Write result to file and stdout
		result := fmt.Sprintf("%s,%s\n", path, hash)
		outputFile.WriteString(result)
		fmt.Print(result)

		return nil
	})

	if err != nil {
		log.Fatal("Error walking through directory:", err)
	}

	fmt.Println("Hashing complete. Results saved to hashes.txt")
}
