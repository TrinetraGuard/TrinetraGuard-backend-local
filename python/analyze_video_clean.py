#!/usr/bin/env python3
import cv2
import json
import os
import sys
import numpy as np
from datetime import timedelta
from collections import defaultdict
import hashlib

def format_timestamp(seconds):
    """Convert seconds to HH:MM:SS format"""
    return str(timedelta(seconds=int(seconds)))

def calculate_face_quality(face_image):
    """Calculate face quality score based on size, clarity, and contrast"""
    # Convert to grayscale
    gray = cv2.cvtColor(face_image, cv2.COLOR_BGR2GRAY)
    
    # Calculate quality metrics
    height, width = gray.shape
    
    # Size quality (prefer larger faces)
    size_score = min(width * height / 10000, 1.0)
    
    # Contrast quality
    contrast = np.std(gray)
    contrast_score = min(contrast / 50, 1.0)
    
    # Sharpness quality (using Laplacian variance)
    laplacian = cv2.Laplacian(gray, cv2.CV_64F)
    sharpness = laplacian.var()
    sharpness_score = min(sharpness / 500, 1.0)
    
    # Overall quality score
    quality_score = (size_score + contrast_score + sharpness_score) / 3
    return quality_score

def are_faces_similar(face1, face2, threshold=0.7):
    """Compare two faces for similarity with higher threshold"""
    # Resize both faces to same size
    face1_resized = cv2.resize(face1, (64, 64))
    face2_resized = cv2.resize(face2, (64, 64))
    
    # Convert to grayscale
    gray1 = cv2.cvtColor(face1_resized, cv2.COLOR_BGR2GRAY)
    gray2 = cv2.cvtColor(face2_resized, cv2.COLOR_BGR2GRAY)
    
    # Calculate correlation coefficient
    correlation = cv2.matchTemplate(gray1, gray2, cv2.TM_CCOEFF_NORMED)[0][0]
    return correlation > threshold

def is_likely_human_face(face_image):
    """Check if the detected face is likely a human face vs logo/object"""
    # Convert to grayscale
    gray = cv2.cvtColor(face_image, cv2.COLOR_BGR2GRAY)
    
    # Check aspect ratio (human faces are roughly 1:1 to 1:1.3)
    height, width = gray.shape
    aspect_ratio = width / height
    if aspect_ratio < 0.7 or aspect_ratio > 1.5:
        return False
    
    # Check for symmetry (human faces are roughly symmetrical)
    left_half = gray[:, :width//2]
    right_half = cv2.flip(gray[:, width//2:], 1)
    if right_half.shape[1] != left_half.shape[1]:
        right_half = right_half[:, :left_half.shape[1]]
    
    if left_half.shape[1] > 0 and right_half.shape[1] > 0:
        symmetry_score = cv2.matchTemplate(left_half, right_half, cv2.TM_CCOEFF_NORMED)[0][0]
        if symmetry_score < 0.3:  # Low symmetry might indicate logo
            return False
    
    # Check for natural skin tone distribution
    hsv = cv2.cvtColor(face_image, cv2.COLOR_BGR2HSV)
    skin_mask = cv2.inRange(hsv, (0, 20, 70), (20, 255, 255))
    skin_ratio = np.sum(skin_mask > 0) / (skin_mask.shape[0] * skin_mask.shape[1])
    
    # If too much or too little skin tone, might not be human
    if skin_ratio < 0.1 or skin_ratio > 0.9:
        return False
    
    return True

def analyze_video(video_path, faces_dir):
    """Analyze video for clear, unique human faces only"""
    
    # Load the pre-trained face detection model
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')
    
    # Open video file
    cap = cv2.VideoCapture(video_path)
    if not cap.isOpened():
        print(json.dumps({"error": "Could not open video file"}))
        return
    
    # Get video properties
    fps = cap.get(cv2.CAP_PROP_FPS)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    duration = total_frames / fps if fps > 0 else 0
    
    print(f"Analyzing video: {video_path}")
    print(f"FPS: {fps}, Total frames: {total_frames}, Duration: {duration:.2f} seconds")
    
    # Store high-quality unique faces
    unique_faces = []
    face_timestamps = defaultdict(list)
    frame_count = 0
    
    # Process every 15th frame to get better quality samples
    sample_rate = 15
    
    while True:
        ret, frame = cap.read()
        if not ret:
            break
            
        frame_count += 1
        
        # Only process every nth frame
        if frame_count % sample_rate != 0:
            continue
            
        # Convert to grayscale for face detection
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        
        # Detect faces with very strict parameters
        faces = face_cascade.detectMultiScale(
            gray,
            scaleFactor=1.3,      # More aggressive scaling
            minNeighbors=10,      # Higher neighbor requirement
            minSize=(100, 100),   # Larger minimum size for clear faces
            maxSize=(400, 400)    # Maximum size to avoid false positives
        )
        
        # Calculate timestamp for this frame
        timestamp_seconds = frame_count / fps
        timestamp_str = format_timestamp(timestamp_seconds)
        
        # Process each face found in this frame
        for (x, y, w, h) in faces:
            # Extract face region
            face_image = frame[y:y+h, x:x+w]
            
            # Skip very small or large faces
            if w < 100 or h < 100 or w > 400 or h > 400:
                continue
            
            # Check if it's likely a human face
            if not is_likely_human_face(face_image):
                continue
            
            # Calculate face quality
            quality_score = calculate_face_quality(face_image)
            
            # Only accept high-quality faces
            if quality_score < 0.6:
                continue
                
            # Check if this face is similar to any existing face
            face_matched = False
            for i, existing_face in enumerate(unique_faces):
                if are_faces_similar(face_image, existing_face['image']):
                    # This is the same person - only update if new face is better quality
                    if quality_score > existing_face['quality']:
                        unique_faces[i] = {
                            'image': face_image,
                            'quality': quality_score,
                            'id': i
                        }
                    face_timestamps[i].append(timestamp_str)
                    face_matched = True
                    break
            
            if not face_matched:
                # This is a new person
                unique_faces.append({
                    'image': face_image,
                    'quality': quality_score,
                    'id': len(unique_faces)
                })
                face_timestamps[len(unique_faces) - 1].append(timestamp_str)
        
        # Progress indicator
        if frame_count % (sample_rate * 30) == 0:
            progress = (frame_count / total_frames) * 100
            print(f"Progress: {progress:.1f}% - Found {len(unique_faces)} unique faces so far")
    
    cap.release()
    
    # Sort faces by quality and take only the best ones
    unique_faces.sort(key=lambda x: x['quality'], reverse=True)
    
    # Limit to maximum 10 unique faces for clarity
    max_faces = min(len(unique_faces), 10)
    unique_faces = unique_faces[:max_faces]
    
    # Save face images and prepare results
    results = {
        "total_people": len(unique_faces),
        "faces": []
    }
    
    # Create faces directory if it doesn't exist
    os.makedirs(faces_dir, exist_ok=True)
    
    # Save each unique face
    for person_id, face_data in enumerate(unique_faces):
        # Save face image
        face_filename = f"person_{person_id + 1}.jpg"
        face_path = os.path.join(faces_dir, face_filename)
        cv2.imwrite(face_path, face_data['image'])
        
        # Add to results
        face_result = {
            "image": face_filename,
            "timestamps": list(set(face_timestamps[face_data['id']]))  # Remove duplicates
        }
        results["faces"].append(face_result)
    
    # Print results as JSON
    print(json.dumps(results, indent=2))

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(json.dumps({"error": "Usage: python3 analyze_video_clean.py <video_path> <faces_dir>"}))
        sys.exit(1)
    
    video_path = sys.argv[1]
    faces_dir = sys.argv[2]
    
    if not os.path.exists(video_path):
        print(json.dumps({"error": f"Video file not found: {video_path}"}))
        sys.exit(1)
    
    analyze_video(video_path, faces_dir) 