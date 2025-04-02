import os
import cv2
import numpy as np
from PIL import Image
import random

# Define paths
base_path = "/home/superclan/Documents/insomnius.tools/examples/perceptualhash-ecommerce"
sample_path = os.path.join(base_path, "sample")
cat_image_path = os.path.join(sample_path, "cat.png")
output_dir = sample_path  # Save in the same directory

def create_modified_versions():
    # Check if the original image exists
    if not os.path.exists(cat_image_path):
        print(f"Error: Source image not found at {cat_image_path}")
        return

    # Create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)

    # Load image with OpenCV and convert from BGR to RGB
    img = cv2.imread(cat_image_path)
    if img is None:
        print(f"Error: Unable to read image at {cat_image_path}")
        return
    img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)  # Convert to RGB

    # 1. Create a strongly blurred version
    blurred_img = cv2.GaussianBlur(img, (121, 121), 0)  # Increased blur strength
    blur_output_path = os.path.join(output_dir, "cat_blurred.png")
    cv2.imwrite(blur_output_path, cv2.cvtColor(blurred_img, cv2.COLOR_RGB2BGR))
    print(f"Created stronger blurred image at {blur_output_path}")

    # 2. Create a resized version
    img_pil = Image.fromarray(img)  # Convert NumPy array to PIL Image
    width, height = img_pil.size
    resized_img = img_pil.resize((width // 80, height // 80))
    resize_output_path = os.path.join(output_dir, "cat_resized.png")
    resized_img.save(resize_output_path)
    print(f"Created resized image at {resize_output_path}")

    # 3. Create a more broken version
    img_array = np.array(img)

    # Add random noise to 30% of pixels
    num_corrupted_pixels = int(0.3 * img_array.size / img_array.shape[-1])
    for _ in range(num_corrupted_pixels):
        y = random.randint(0, img_array.shape[0] - 1)  # Height
        x = random.randint(0, img_array.shape[1] - 1)  # Width
        if random.random() < 0.5:
            img_array[y, x] = [0, 0, 0]  # Pure black pixels
        else:
            img_array[y, x] = [255, 255, 255]  # Pure white pixels

    # Add more corruption blocks
    for _ in range(50):  # Increased number of blocks
        block_w = random.randint(10, 30)
        block_h = random.randint(10, 30)
        block_x = random.randint(0, img_array.shape[1] - block_w)
        block_y = random.randint(0, img_array.shape[0] - block_h)
        img_array[block_y:block_y + block_h, block_x:block_x + block_w] = np.random.randint(0, 256, size=(block_h, block_w, 3), dtype=np.uint8)

    # Apply a stronger motion blur to random areas
    if random.random() < 0.8:  # Apply in 80% of cases
        kernel_size = random.choice([9, 11, 15])  # Stronger blur
        motion_blur_kernel = np.zeros((kernel_size, kernel_size))
        motion_blur_kernel[:, kernel_size // 2] = 1
        motion_blur_kernel /= kernel_size
        img_array = cv2.filter2D(img_array, -1, motion_blur_kernel)

    # Save the broken image with aggressive JPEG compression (10% quality)
    broken_img = Image.fromarray(img_array)
    broken_output_path = os.path.join(output_dir, "cat_broken.jpg")
    broken_img.save(broken_output_path, quality=10)  # Low quality for more artifacts
    print(f"Created heavily broken image at {broken_output_path}")

if __name__ == "__main__":
    create_modified_versions()
    print("All modified image versions have been created successfully!")
