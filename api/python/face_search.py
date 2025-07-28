import sys
import json
import os
import argparse
import cv2
import face_recognition
import numpy as np
from PIL import Image
import warnings

# Suppress all warnings to ensure clean JSON output
warnings.filterwarnings("ignore")

def load_and_encode_image(image_path):
    """Load and encode a face image"""
    try:
        # Load image
        image = face_recognition.load_image_file(image_path)
        
        # Find face locations
        face_locations = face_recognition.face_locations(image)
        
        if len(face_locations) == 0:
            return None
        
        # Get face encodings
        face_encodings = face_recognition.face_encodings(image, face_locations)
        
        if len(face_encodings) == 0:
            return None
        
        # Return the first face encoding
        return face_encodings[0]
        
    except Exception as e:
        print(f"Error loading image {image_path}: {str(e)}")
        return None

def compare_faces(search_encoding, face_images, similarity_threshold=0.6):
    """Compare search face with stored face images"""
    matched_faces = []
    
    for face_image in face_images:
        try:
            # Remove 'faces/' prefix if present and construct full path
            clean_face_image = face_image.replace('faces/', '')
            face_path = f"../storage/faces/{clean_face_image}"
            
            print(f"Checking face image: {face_path}")
            
            if not os.path.exists(face_path):
                print(f"Face image not found: {face_path}")
                continue
            
            # Load and encode the stored face
            stored_encoding = load_and_encode_image(face_path)
            
            if stored_encoding is None:
                continue
            
            # Compare faces
            distance = face_recognition.face_distance([stored_encoding], search_encoding)[0]
            similarity = 1 - distance
            
            # If similarity is above threshold, consider it a match
            if similarity >= similarity_threshold:
                matched_faces.append(face_image)  # Keep original path for response
                print(f"Match found: {face_image} (similarity: {similarity:.3f})")
            
        except Exception as e:
            print(f"Error comparing with {face_image}: {str(e)}")
            continue
    
    return matched_faces

def main():
    parser = argparse.ArgumentParser(description="Search for faces in stored images")
    parser.add_argument("search_image", help="Path to the search image")
    parser.add_argument("--face-images", help="Comma-separated list of face images to compare")
    parser.add_argument("--threshold", type=float, default=0.6, help="Similarity threshold (default: 0.6)")
    
    args = parser.parse_args()
    
    if not os.path.exists(args.search_image):
        print(json.dumps({"error": "Search image not found"}))
        sys.exit(1)
    
    try:
        # Load and encode the search image
        print(f"Loading search image: {args.search_image}")
        search_encoding = load_and_encode_image(args.search_image)
        
        if search_encoding is None:
            print(json.dumps({"error": "No face found in search image"}))
            sys.exit(1)
        
        # Parse face images list
        if not args.face_images:
            print(json.dumps({"error": "No face images provided"}))
            sys.exit(1)
        
        face_images = [img.strip() for img in args.face_images.split(",") if img.strip()]
        
        if not face_images:
            print(json.dumps({"error": "No valid face images provided"}))
            sys.exit(1)
        
        # Compare faces
        matched_faces = compare_faces(search_encoding, face_images, args.threshold)
        
        # Prepare result
        result = {
            "matched_faces": matched_faces,
            "total_faces_checked": len(face_images),
            "matches_found": len(matched_faces)
        }
        
        sys.stdout.flush()  # Clear any buffered output
        print(json.dumps(result, indent=2))
        sys.stdout.flush()  # Ensure output is sent
        
    except Exception as e:
        error_response = {
            "error": f"Face search failed: {str(e)}",
            "matched_faces": [],
            "total_faces_checked": 0,
            "matches_found": 0
        }
        sys.stdout.flush()  # Clear any buffered output
        print(json.dumps(error_response, indent=2))
        sys.stdout.flush()  # Ensure output is sent
        sys.exit(1)

if __name__ == "__main__":
    main() 