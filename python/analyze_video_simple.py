#!/usr/bin/env python3
import cv2
import json
import os
import sys
import numpy as np
from datetime import timedelta
from collections import defaultdict

def format_timestamp(seconds):
    """Convert seconds to HH:MM:SS format"""
    return str(timedelta(seconds=int(seconds)))

def analyze_video(video_path, faces_dir):
    """Analyze video for faces using OpenCV's Haar Cascade"""
    
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
    
    # Store detected faces
    face_timestamps = defaultdict(list)
    frame_count = 0
    
    # Process every 10th frame to speed up analysis
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
        
        # Detect faces
        faces = face_cascade.detectMultiScale(
            gray,
            scaleFactor=1.1,
            minNeighbors=5,
            minSize=(30, 30)
        )
        
        # Calculate timestamp for this frame
        timestamp_seconds = frame_count / fps
        timestamp_str = format_timestamp(timestamp_seconds)
        
        # Process each face found in this frame
        for i, (x, y, w, h) in enumerate(faces):
            # For simplicity, we'll treat each detection as a unique person
            # In a real system, you'd want to track faces across frames
            person_id = len(face_timestamps)
            face_timestamps[person_id].append(timestamp_str)
        
        # Progress indicator
        if frame_count % (sample_rate * 30) == 0:
            progress = (frame_count / total_frames) * 100
            print(f"Progress: {progress:.1f}% - Found {len(face_timestamps)} unique faces so far")
    
    cap.release()
    
    # Save face images and prepare results
    results = {
        "total_people": len(face_timestamps),
        "faces": []
    }
    
    # Create faces directory if it doesn't exist
    os.makedirs(faces_dir, exist_ok=True)
    
    # Reopen video to extract face images
    cap = cv2.VideoCapture(video_path)
    
    for person_id, timestamps in face_timestamps.items():
        # Find a frame where this person appears
        cap.set(cv2.CAP_PROP_POS_FRAMES, 0)
        best_face_image = None
        
        while True:
            ret, frame = cap.read()
            if not ret:
                break
                
            gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
            faces = face_cascade.detectMultiScale(
                gray,
                scaleFactor=1.1,
                minNeighbors=5,
                minSize=(30, 30)
            )
            
            if len(faces) > 0:
                # Extract the first face found
                x, y, w, h = faces[0]
                face_image = frame[y:y+h, x:x+w]
                
                # Resize for consistency
                face_image = cv2.resize(face_image, (200, 200))
                best_face_image = face_image
                break
        
        if best_face_image is not None:
            # Save face image
            face_filename = f"person_{person_id + 1}.jpg"
            face_path = os.path.join(faces_dir, face_filename)
            cv2.imwrite(face_path, best_face_image)
            
            # Add to results
            face_result = {
                "image": face_filename,
                "timestamps": list(set(timestamps))  # Remove duplicates
            }
            results["faces"].append(face_result)
    
    cap.release()
    
    # Print results as JSON
    print(json.dumps(results, indent=2))

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(json.dumps({"error": "Usage: python3 analyze_video_simple.py <video_path> <faces_dir>"}))
        sys.exit(1)
    
    video_path = sys.argv[1]
    faces_dir = sys.argv[2]
    
    if not os.path.exists(video_path):
        print(json.dumps({"error": f"Video file not found: {video_path}"}))
        sys.exit(1)
    
    analyze_video(video_path, faces_dir) 