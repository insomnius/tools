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
	type ImageHash struct {
		Path string
		Hash string
	}

	var imageHashes []ImageHash

	err := filepath.Walk("./sample", func(path string, info os.FileInfo, err error) error {
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
		debugFolder := "./debug"
		if _, err := os.Stat(debugFolder); os.IsNotExist(err) {
			os.Mkdir(debugFolder, 0755)
		}

		// Configure perceptual hash with debug options
		conf := perceptualhash.Config{
			Debug: true,
		}

		if conf.Debug {
			baseName := filepath.Base(path)
			targetFolder := filepath.Join(debugFolder, baseName)
			if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
				os.Mkdir(targetFolder, 0755)
			}
			conf.DebugParameter.PreprocessedImagePath = filepath.Join(targetFolder, "preprocessed_"+baseName)
			conf.DebugParameter.VisualizedImagePath = filepath.Join(targetFolder, "visualized_"+baseName)
		}

		// Generate hash
		hash, err := perceptualhash.FromPath(path, conf)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", path, err)
			return nil // Continue with next file
		}

		// Write result to file and stdout
		result := fmt.Sprintf("%s,%s\n", path, hash)
		fmt.Print(result)

		imageHashes = append(imageHashes, ImageHash{
			Path: path,
			Hash: hash,
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	var originalImageHash ImageHash
	for _, ih := range imageHashes {
		if ih.Path == "sample/cat.png" {
			originalImageHash = ih
			break
		}
	}

	fmt.Printf("Original Image Hash: %s\n", originalImageHash.Hash)
	fmt.Printf("Original Image Path: %s\n", originalImageHash.Path)
	fmt.Println("-------------------------")
	for _, ih := range imageHashes {
		fmt.Printf("Path: %s, Hash: %s\n", ih.Path, ih.Hash)
		if ih.Path == "sample/cat.png" {
			continue
		}
		distance, err := perceptualhash.CompareHashes(originalImageHash.Hash, ih.Hash)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Distance between %s and %s: %d\n", originalImageHash.Path, ih.Path, distance)
		if distance < 10 {
			fmt.Printf("Similar image found: %s\n", ih.Path)
		}
		fmt.Println("-------------------------")
	}

}
