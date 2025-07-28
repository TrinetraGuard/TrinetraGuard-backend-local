#!/usr/bin/env python3
import cv2
import numpy as np
import json
import os
import sys
from datetime import timedelta
import math

def calculate_face_similarity(face1, face2):
    """Calculate similarity between two face images using multiple methods"""
    try:
        # Resize both faces to same size for comparison
        size = (100, 100)
        face1_resized = cv2.resize(face1, size)
        face2_resized = cv2.resize(face2, size)
        
        # Convert to grayscale
        gray1 = cv2.cvtColor(face1_resized, cv2.COLOR_BGR2GRAY)
        gray2 = cv2.cvtColor(face2_resized, cv2.COLOR_BGR2GRAY)
        
        # Method 1: Template matching
        result = cv2.matchTemplate(gray1, gray2, cv2.TM_CCOEFF_NORMED)
        template_similarity = result[0][0]
        
        # Method 2: Histogram comparison
        hist1 = cv2.calcHist([gray1], [0], None, [256], [0, 256])
        hist2 = cv2.calcHist([gray2], [0], None, [256], [0, 256])
        hist_similarity = cv2.compareHist(hist1, hist2, cv2.HISTCMP_CORREL)
        
        # Method 3: Structural similarity (SSIM-like)
        diff = cv2.absdiff(gray1, gray2)
        ssim_similarity = 1.0 - (np.mean(diff) / 255.0)
        
        # Combine all similarities with weights
        combined_similarity = (
            template_similarity * 0.5 +
            hist_similarity * 0.3 +
            ssim_similarity * 0.2
        )
        
        return combined_similarity
    except:
        return 0.0

def calculate_face_quality(face_img):
    """Calculate quality score for a face image"""
    try:
        # Convert to grayscale
        gray = cv2.cvtColor(face_img, cv2.COLOR_BGR2GRAY)
        
        # Calculate sharpness using Laplacian variance
        laplacian_var = cv2.Laplacian(gray, cv2.CV_64F).var()
        
        # Calculate contrast
        contrast = gray.std()
        
        # Calculate brightness
        brightness = gray.mean()
        
        # Calculate size score (prefer larger faces)
        height, width = face_img.shape[:2]
        size_score = min(height * width / 25000, 1.0)  # Even higher size requirement
        
        # Calculate aspect ratio (prefer faces closer to 1:1)
        aspect_ratio = width / height
        aspect_score = 1.0 - abs(aspect_ratio - 1.0) * 0.5
        
        # Calculate symmetry (basic check)
        left_half = gray[:, :width//2]
        right_half = cv2.flip(gray[:, width//2:], 1)
        if left_half.shape == right_half.shape:
            symmetry = 1.0 - np.mean(np.abs(left_half.astype(float) - right_half.astype(float))) / 255
        else:
            symmetry = 0.5
        
        # Combine all scores with much higher weight on size and quality
        quality_score = (
            laplacian_var * 0.15 +     # Sharpness
            contrast * 0.15 +           # Contrast
            (255 - abs(brightness - 128)) * 0.1 +  # Brightness
            size_score * 0.45 +         # Size (much higher weight)
            aspect_score * 0.1 +        # Aspect ratio
            symmetry * 0.05             # Symmetry
        )
        
        return quality_score
    except:
        return 0.0

def is_complete_human_face(face_img):
    """Check if the detected face is a complete human face (not just parts)"""
    try:
        height, width = face_img.shape[:2]
        
        # Check aspect ratio (human faces are roughly 1:1 to 1:1.3)
        aspect_ratio = width / height
        if aspect_ratio < 0.9 or aspect_ratio > 1.4:
            return False
        
        # Check size (must be large enough to be a complete face)
        if width < 180 or height < 180:
            return False
        
        # Check if face is too large (likely a false positive)
        if width > 350 or height > 350:
            return False
        
        # Convert to HSV and check skin tone distribution
        hsv = cv2.cvtColor(face_img, cv2.COLOR_BGR2HSV)
        
        # Define skin tone ranges (very restrictive)
        lower_skin = np.array([0, 40, 90], dtype=np.uint8)
        upper_skin = np.array([20, 255, 255], dtype=np.uint8)
        
        # Create mask for skin tones
        skin_mask = cv2.inRange(hsv, lower_skin, upper_skin)
        skin_pixels = cv2.countNonZero(skin_mask)
        total_pixels = width * height
        skin_ratio = skin_pixels / total_pixels
        
        # Must have reasonable skin tone ratio (not too much, not too little)
        if skin_ratio < 0.25 or skin_ratio > 0.65:
            return False
        
        # Check for face-like features using edge detection
        gray = cv2.cvtColor(face_img, cv2.COLOR_BGR2GRAY)
        edges = cv2.Canny(gray, 50, 150)
        edge_density = np.sum(edges > 0) / (width * height)
        
        # Too many edges might indicate noise or partial face
        if edge_density > 0.2:
            return False
        
        # Check for face symmetry more strictly
        left_half = gray[:, :width//2]
        right_half = cv2.flip(gray[:, width//2:], 1)
        if left_half.shape == right_half.shape:
            symmetry_score = 1.0 - np.mean(np.abs(left_half.astype(float) - right_half.astype(float))) / 255
            if symmetry_score < 0.45:  # Must be reasonably symmetrical
                return False
        
        # Additional check: look for eyes using eye cascade
        eye_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_eye.xml')
        eyes = eye_cascade.detectMultiScale(gray, 1.1, 3)
        
        # Must have at least 1 eye detected to be a complete face
        if len(eyes) < 1:
            return False
        
        return True
        
    except:
        return False

def analyze_video(video_path, faces_dir):
    """Analyze video and detect only complete, unique human faces"""
    
    # Create faces directory if it doesn't exist
    os.makedirs(faces_dir, exist_ok=True)
    
    # Load the face cascade classifier
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')
    
    # Open the video
    cap = cv2.VideoCapture(video_path)
    if not cap.isOpened():
        print(f"Error: Could not open video {video_path}")
        return None
    
    # Get video properties
    fps = cap.get(cv2.CAP_PROP_FPS)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    duration = total_frames / fps if fps > 0 else 0
    
    print(f"Analyzing video: {video_path}")
    print(f"FPS: {fps}, Total frames: {total_frames}, Duration: {duration:.2f} seconds")
    
    # Store detected faces
    detected_faces = []
    frame_count = 0
    
    # Process every 10th frame for better quality samples
    frame_skip = 10
    
    while True:
        ret, frame = cap.read()
        if not ret:
            break
        
        frame_count += 1
        
        # Process every nth frame to get better quality samples
        if frame_count % frame_skip != 0:
            continue
        
        # Convert to grayscale for face detection
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        
        # Detect faces with very strict parameters for complete faces only
        faces = face_cascade.detectMultiScale(
            gray,
            scaleFactor=1.03,     # Very precise scaling
            minNeighbors=15,      # Very high threshold for better detection
            minSize=(180, 180),   # Large minimum size for complete faces
            maxSize=(350, 350)    # Reasonable maximum size
        )
        
        # Process each detected face
        for (x, y, w, h) in faces:
            # Extract face region
            face_img = frame[y:y+h, x:x+w]
            
            # Check if it's a complete human face
            if not is_complete_human_face(face_img):
                continue
            
            # Calculate quality score
            quality_score = calculate_face_quality(face_img)
            
            # Only keep very high-quality faces
            if quality_score < 40:  # Very high threshold for better quality
                continue
            
            # Calculate timestamp
            timestamp_seconds = frame_count / fps
            timestamp = str(timedelta(seconds=int(timestamp_seconds)))
            
            # Check if this face is similar to any existing face
            is_duplicate = False
            best_match_index = -1
            best_similarity = 0.0
            
            for i, existing_face in enumerate(detected_faces):
                # Load existing face image for comparison
                existing_img_path = os.path.join(faces_dir, existing_face['image'])
                if os.path.exists(existing_img_path):
                    existing_img = cv2.imread(existing_img_path)
                    if existing_img is not None:
                        similarity = calculate_face_similarity(face_img, existing_img)
                        
                        if similarity > best_similarity:
                            best_similarity = similarity
                            best_match_index = i
                        
                        # If similarity is high, it's a duplicate
                        if similarity > 0.9:  # Very high threshold for strict grouping
                            is_duplicate = True
                            break
            
            if is_duplicate:
                # Update existing face with better quality image if available
                if best_match_index >= 0 and quality_score > existing_face.get('quality', 0):
                    # Save better quality image
                    face_filename = f"person_{best_match_index + 1}.jpg"
                    face_path = os.path.join(faces_dir, face_filename)
                    cv2.imwrite(face_path, face_img)
                    
                    # Update existing face data
                    detected_faces[best_match_index]['quality'] = quality_score
                    detected_faces[best_match_index]['image'] = face_filename
                
                # Add timestamp to existing face
                if best_match_index >= 0:
                    if timestamp not in detected_faces[best_match_index]['timestamps']:
                        detected_faces[best_match_index]['timestamps'].append(timestamp)
            else:
                # New face detected
                face_index = len(detected_faces) + 1
                face_filename = f"person_{face_index}.jpg"
                face_path = os.path.join(faces_dir, face_filename)
                
                # Save face image
                cv2.imwrite(face_path, face_img)
                
                # Add to detected faces
                detected_faces.append({
                    'image': face_filename,
                    'timestamps': [timestamp],
                    'quality': quality_score
                })
        
        # Progress update
        if frame_count % (total_frames // 10) == 0:
            progress = (frame_count / total_frames) * 100
            print(f"Progress: {progress:.1f}% - Found {len(detected_faces)} unique faces so far")
    
    # Release video capture
    cap.release()
    
    # Sort faces by quality and limit to top results
    detected_faces.sort(key=lambda x: x.get('quality', 0), reverse=True)
    
    # Limit to maximum 8 unique faces for clarity
    max_faces = min(len(detected_faces), 8)
    detected_faces = detected_faces[:max_faces]
    
    # Prepare result
    result = {
        "total_people": len(detected_faces),
        "faces": []
    }
    
    for i, face in enumerate(detected_faces):
        result["faces"].append({
            "image": face['image'],
            "timestamps": sorted(face['timestamps'])
        })
    
    # Print result as JSON
    print(json.dumps(result, indent=2))
    
    return result

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 analyze_video_clean.py <video_path> <faces_dir>")
        sys.exit(1)
    
    video_path = sys.argv[1]
    faces_dir = sys.argv[2]
    
    result = analyze_video(video_path, faces_dir)
    if result is None:
        sys.exit(1) 