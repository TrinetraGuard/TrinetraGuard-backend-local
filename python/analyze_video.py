#!/usr/bin/env python3
import cv2
import face_recognition
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
    """Analyze video for faces and return analysis results"""
    
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
    face_encodings = []
    face_locations = []
    face_timestamps = defaultdict(list)
    frame_count = 0
    
    # Process every 10th frame to speed up analysis (adjust as needed)
    sample_rate = 10
    
    while True:
        ret, frame = cap.read()
        if not ret:
            break
            
        frame_count += 1
        
        # Only process every nth frame
        if frame_count % sample_rate != 0:
            continue
            
        # Convert BGR to RGB (face_recognition uses RGB)
        rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        
        # Find faces in the frame
        face_locations_in_frame = face_recognition.face_locations(rgb_frame)
        face_encodings_in_frame = face_recognition.face_encodings(rgb_frame, face_locations_in_frame)
        
        # Calculate timestamp for this frame
        timestamp_seconds = frame_count / fps
        timestamp_str = format_timestamp(timestamp_seconds)
        
        # Process each face found in this frame
        for i, (face_encoding, face_location) in enumerate(zip(face_encodings_in_frame, face_locations_in_frame)):
            # Check if this face matches any previously seen face
            face_matched = False
            for j, existing_encoding in enumerate(face_encodings):
                # Compare face encodings
                if face_recognition.compare_faces([existing_encoding], face_encoding, tolerance=0.6)[0]:
                    # This is the same person
                    face_timestamps[j].append(timestamp_str)
                    face_matched = True
                    break
            
            if not face_matched:
                # This is a new person
                face_encodings.append(face_encoding)
                face_locations.append(face_location)
                face_timestamps[len(face_encodings) - 1].append(timestamp_str)
        
        # Progress indicator
        if frame_count % (sample_rate * 30) == 0:
            progress = (frame_count / total_frames) * 100
            print(f"Progress: {progress:.1f}% - Found {len(face_encodings)} unique faces so far")
    
    cap.release()
    
    # Save face images and prepare results
    results = {
        "total_people": len(face_encodings),
        "faces": []
    }
    
    # Create faces directory if it doesn't exist
    os.makedirs(faces_dir, exist_ok=True)
    
    # Reopen video to extract face images
    cap = cv2.VideoCapture(video_path)
    
    for person_id, (face_encoding, face_location) in enumerate(zip(face_encodings, face_locations)):
        # Find a frame where this person appears clearly
        cap.set(cv2.CAP_PROP_POS_FRAMES, 0)
        best_face_image = None
        
        while True:
            ret, frame = cap.read()
            if not ret:
                break
                
            rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            face_locations_in_frame = face_recognition.face_locations(rgb_frame)
            face_encodings_in_frame = face_recognition.face_encodings(rgb_frame, face_locations_in_frame)
            
            # Check if this person appears in this frame
            for face_encoding_in_frame in face_encodings_in_frame:
                if face_recognition.compare_faces([face_encoding], face_encoding_in_frame, tolerance=0.6)[0]:
                    # Extract face region
                    top, right, bottom, left = face_recognition.face_locations(rgb_frame)[0]
                    face_image = rgb_frame[top:bottom, left:right]
                    
                    # Resize for consistency
                    face_image = cv2.resize(face_image, (200, 200))
                    best_face_image = face_image
                    break
            
            if best_face_image is not None:
                break
        
        # Save face image
        face_filename = f"person_{person_id + 1}.jpg"
        face_path = os.path.join(faces_dir, face_filename)
        cv2.imwrite(face_path, cv2.cvtColor(best_face_image, cv2.COLOR_RGB2BGR))
        
        # Add to results
        face_result = {
            "image": face_filename,
            "timestamps": list(set(face_timestamps[person_id]))  # Remove duplicates
        }
        results["faces"].append(face_result)
    
    cap.release()
    
    # Print results as JSON
    print(json.dumps(results, indent=2))

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(json.dumps({"error": "Usage: python3 analyze_video.py <video_path> <faces_dir>"}))
        sys.exit(1)
    
    video_path = sys.argv[1]
    faces_dir = sys.argv[2]
    
    if not os.path.exists(video_path):
        print(json.dumps({"error": f"Video file not found: {video_path}"}))
        sys.exit(1)
    
    analyze_video(video_path, faces_dir) 