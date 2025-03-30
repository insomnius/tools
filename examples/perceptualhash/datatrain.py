import fiftyone as fo
import fiftyone.zoo as foz
import os
import shutil

# Download a dataset (COCO-2017 as an example)
dataset = foz.load_zoo_dataset("coco-2017", split="validation")

# Choose where to save the images
output_dir = "images"
if os.path.exists(output_dir):
    shutil.rmtree(output_dir)  # Clear old files
os.makedirs(output_dir, exist_ok=True)

# Export images
for sample in dataset:
    image_path = sample.filepath
    shutil.copy(image_path, output_dir)

print(f"Downloaded {len(dataset)} images to {output_dir}")
