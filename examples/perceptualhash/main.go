package main

import (
	"fmt"
	"log"

	"github.com/insomnius/tools/perceptualhash"
)

func main() {
	conf := perceptualhash.Config{
		Debug: true,
	}
	conf.DebugParameter.PreprocessedImagePath = "./sample/sample-image-preprocessed.png"
	conf.DebugParameter.VisualizedImagePath = "./sample/sample-image-visualized.png"

	hash, err := perceptualhash.FromPath("./sample/sample-image.png", conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hash", hash)
}
