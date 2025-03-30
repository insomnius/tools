package main

import (
	"fmt"
	"log"
	"math/bits"
	"os"
	"path/filepath"
	"strings"

	"github.com/insomnius/tools/perceptualhash"
)

type ImageHash struct {
	Path string
	Hash string
}

func hammingDistance(hash1, hash2 string) int {
	// Skip if lengths don't match
	if len(hash1) != len(hash2) {
		return -1
	}

	// Count differing bits
	distance := 0
	for i := 0; i < len(hash1); i++ {
		// XOR the bytes and count the bits
		distance += bits.OnesCount8(hash1[i] ^ hash2[i])
	}
	return distance
}

func main() {
	// Path to the folder containing images
	imagesFolder := "./images"

	// Output file to save the results
	outputFile, err := os.Create("hashes.txt")
	if err != nil {
		log.Fatal("Failed to create output file:", err)
	}
	defer outputFile.Close()

	// Slice to store all image hashes
	var imageHashes []ImageHash

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
		outputFile.WriteString(result)
		fmt.Print(result)

		// Store hash for confusion matrix calculation
		imageHashes = append(imageHashes, ImageHash{
			Path: path,
			Hash: hash,
		})

		return nil
	})

	if err != nil {
		log.Fatal("Error walking through directory:", err)
	}

	fmt.Println("Hashing complete. Results saved to hashes.txt")

	// Calculate confusion matrix
	calculateConfusionMatrix(imageHashes)
}

func calculateConfusionMatrix(imageHashes []ImageHash) {
	// Threshold for considering two images similar based on Hamming distance
	const similarityThreshold = 10

	// Confusion matrix values
	var truePositives, falsePositives, trueNegatives, falseNegatives int

	// Compare each pair of images
	for i := 0; i < len(imageHashes); i++ {
		for j := i + 1; j < len(imageHashes); j++ {
			hash1 := imageHashes[i].Hash
			hash2 := imageHashes[j].Hash
			path1 := imageHashes[i].Path
			path2 := imageHashes[j].Path

			// Compute Hamming distance between hashes
			distance := hammingDistance(hash1, hash2)
			if distance < 0 {
				fmt.Printf("Cannot compare hashes of different lengths: %s and %s\n", path1, path2)
				continue
			}

			// Determine if the images should be similar based on their names
			// This is a simple heuristic; adjust according to your dataset
			// Assuming images with the same prefix (before first underscore) are similar
			base1 := filepath.Base(path1)
			base2 := filepath.Base(path2)
			prefix1 := strings.Split(base1, "_")[0]
			prefix2 := strings.Split(base2, "_")[0]
			shouldBeSimilar := prefix1 == prefix2

			// Determine if hashes indicate similarity
			hashIndicatesSimilar := distance <= similarityThreshold

			// Update confusion matrix
			if shouldBeSimilar && hashIndicatesSimilar {
				truePositives++
			} else if !shouldBeSimilar && hashIndicatesSimilar {
				falsePositives++
			} else if !shouldBeSimilar && !hashIndicatesSimilar {
				trueNegatives++
			} else if shouldBeSimilar && !hashIndicatesSimilar {
				falseNegatives++
			}
		}
	}

	// Calculate metrics
	accuracy := float64(truePositives+trueNegatives) / float64(truePositives+trueNegatives+falsePositives+falseNegatives)

	var precision, recall, f1Score float64

	if truePositives+falsePositives > 0 {
		precision = float64(truePositives) / float64(truePositives+falsePositives)
	}

	if truePositives+falseNegatives > 0 {
		recall = float64(truePositives) / float64(truePositives+falseNegatives)
	}

	if precision+recall > 0 {
		f1Score = 2 * precision * recall / (precision + recall)
	}

	// Output confusion matrix
	fmt.Println("\nConfusion Matrix:")
	fmt.Println("=================")
	fmt.Printf("True Positives: %d\n", truePositives)
	fmt.Printf("False Positives: %d\n", falsePositives)
	fmt.Printf("True Negatives: %d\n", trueNegatives)
	fmt.Printf("False Negatives: %d\n", falseNegatives)
	fmt.Println("\nMetrics:")
	fmt.Printf("Accuracy: %.4f\n", accuracy)
	fmt.Printf("Precision: %.4f\n", precision)
	fmt.Printf("Recall: %.4f\n", recall)
	fmt.Printf("F1 Score: %.4f\n", f1Score)
}
