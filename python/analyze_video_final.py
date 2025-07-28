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

def calculate_face_hash(face_image):
    """Calculate a simple hash for face comparison"""
    # Convert to grayscale and resize for consistent hashing
    gray = cv2.cvtColor(face_image, cv2.COLOR_BGR2GRAY)
    resized = cv2.resize(gray, (16, 16))  # Smaller size for more aggressive grouping
    # Calculate hash based on pixel values
    hash_value = hashlib.md5(resized.tobytes()).hexdigest()
    return hash_value

def are_faces_similar(face1, face2, threshold=0.6):  # Lower threshold for more aggressive grouping
    """Compare two faces for similarity using structural similarity"""
    # Resize both faces to same size
    face1_resized = cv2.resize(face1, (50, 50))  # Smaller size
    face2_resized = cv2.resize(face2, (50, 50))
    
    # Convert to grayscale
    gray1 = cv2.cvtColor(face1_resized, cv2.COLOR_BGR2GRAY)
    gray2 = cv2.cvtColor(face2_resized, cv2.COLOR_BGR2GRAY)
    
    # Calculate correlation coefficient
    correlation = cv2.matchTemplate(gray1, gray2, cv2.TM_CCOEFF_NORMED)[0][0]
    return correlation > threshold

def analyze_video(video_path, faces_dir):
    """Analyze video for faces using OpenCV's Haar Cascade with aggressive grouping"""
    
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
    
    # Store detected faces with aggressive grouping
    unique_faces = []
    face_timestamps = defaultdict(list)
    frame_count = 0
    
    # Process every 10th frame to reduce processing time
    sample_rate = 10
    
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
        
        # Detect faces with stricter parameters
        faces = face_cascade.detectMultiScale(
            gray,
            scaleFactor=1.2,  # More aggressive scaling
            minNeighbors=8,   # Higher neighbor requirement
            minSize=(80, 80)  # Larger minimum size for better quality
        )
        
        # Calculate timestamp for this frame
        timestamp_seconds = frame_count / fps
        timestamp_str = format_timestamp(timestamp_seconds)
        
        # Process each face found in this frame
        for (x, y, w, h) in faces:
            # Extract face region
            face_image = frame[y:y+h, x:x+w]
            
            # Skip very small faces
            if w < 80 or h < 80:
                continue
                
            # Check if this face is similar to any existing face
            face_matched = False
            for i, existing_face in enumerate(unique_faces):
                if are_faces_similar(face_image, existing_face['image']):
                    # This is the same person
                    face_timestamps[i].append(timestamp_str)
                    face_matched = True
                    break
            
            if not face_matched:
                # This is a new person
                face_hash = calculate_face_hash(face_image)
                unique_faces.append({
                    'image': face_image,
                    'hash': face_hash,
                    'id': len(unique_faces)
                })
                face_timestamps[len(unique_faces) - 1].append(timestamp_str)
        
        # Progress indicator
        if frame_count % (sample_rate * 30) == 0:
            progress = (frame_count / total_frames) * 100
            print(f"Progress: {progress:.1f}% - Found {len(unique_faces)} unique faces so far")
    
    cap.release()
    
    # Limit to maximum 20 unique faces to avoid overwhelming results
    max_faces = min(len(unique_faces), 20)
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
            "timestamps": list(set(face_timestamps[person_id]))  # Remove duplicates
        }
        results["faces"].append(face_result)
    
    # Print results as JSON
    print(json.dumps(results, indent=2))

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(json.dumps({"error": "Usage: python3 analyze_video_final.py <video_path> <faces_dir>"}))
        sys.exit(1)
    
    video_path = sys.argv[1]
    faces_dir = sys.argv[2]
    
    if not os.path.exists(video_path):
        print(json.dumps({"error": f"Video file not found: {video_path}"}))
        sys.exit(1)
    
    analyze_video(video_path, faces_dir) 