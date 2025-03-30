# Repository Overview

This repository contains tools and utilities for Insomnius, designed to simplify and enhance various workflows.

## Packages

### 1. Perceptual Hash (`perceptualhash`)
A package for generating perceptual hashes from images. It includes:
- Image preprocessing.
- Hash generation using Discrete Cosine Transform (DCT).
- Debugging tools for visualizing the hash.

#### Example Usage
Refer to the `examples/perceptualhash` folder for:
- `main.go`: Demonstrates how to generate a perceptual hash for an image.
- `datatrain.py`: A Python script for downloading and exporting datasets like COCO-2017.

## Usage

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/insomnius.tools.git
   ```

2. Install dependencies:
   - For Go packages:
     ```bash
     go mod tidy
     ```

3. Run the tools as needed.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This repository is licensed under MIT license. See the `LICENSE` file for details.
