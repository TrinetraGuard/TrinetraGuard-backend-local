#!/usr/bin/env python3
"""
Video Face Detection and Recognition System
High-accuracy face detection and deduplication using face_recognition library
"""

import sys
import json
import os
import argparse
import time
from pathlib import Path
import cv2
import face_recognition
import numpy as np
from PIL import Image
import warnings

# Suppress all warnings to ensure clean JSON output
warnings.filterwarnings("ignore")

class FaceProcessor:
    def __init__(self, video_path, video_id=None, fps=1, threshold=0.6):
        self.video_path = video_path
        self.fps = fps
        self.threshold = threshold
        self.known_faces = []
        self.known_encodings = []
        self.face_count = 0
        
        # Create faces directory if it doesn't exist
        faces_dir = Path("../storage/faces")
        faces_dir.mkdir(exist_ok=True)
        
        # Use provided video ID or generate from filename
        if video_id:
            self.video_id = video_id
        else:
            video_filename = Path(video_path).stem
            self.video_id = video_filename
        
    def extract_frames(self, video_path):
        """Extract frames from video at specified FPS"""
        cap = cv2.VideoCapture(video_path)
        if not cap.isOpened():
            raise ValueError("Could not open video file")
            
        # Get video properties
        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        video_fps = cap.get(cv2.CAP_PROP_FPS)
        duration = total_frames / video_fps
        
        print(f"Video info: {total_frames} frames, {video_fps:.2f} fps, {duration:.2f}s duration")
        
        frames = []
        frame_interval = int(video_fps / self.fps)
        
        frame_count = 0
        while True:
            ret, frame = cap.read()
            if not ret:
                break
                
            if frame_count % frame_interval == 0:
                # Convert BGR to RGB
                rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
                frames.append(rgb_frame)
                
            frame_count += 1
            
        cap.release()
        print(f"Extracted {len(frames)} frames at {self.fps} fps")
        return frames
        
    def process_faces(self, frame, frame_num):
        """Process faces in a single frame"""
        # Find face locations
        face_locations = face_recognition.face_locations(frame)
        face_encodings = face_recognition.face_encodings(frame, face_locations)
        
        print(f"Found {len(face_locations)} faces in frame {frame_num}")
        
        new_faces = []
        
        for i, (face_location, face_encoding) in enumerate(zip(face_locations, face_encodings)):
            # Check if this face is similar to any known face
            if len(self.known_encodings) > 0:
                matches = face_recognition.compare_faces(self.known_encodings, face_encoding, tolerance=self.threshold)
                if any(matches):
                    print("Duplicate face detected (skipping)")
                    continue
            
            # This is a new face
            self.face_count += 1
            print(f"New face detected! Face #{self.face_count}")
            
            # Save the face image with unique filename
            top, right, bottom, left = face_location
            face_image = frame[top:bottom, left:right]
            
            # Convert to PIL Image and save with unique name
            pil_image = Image.fromarray(face_image)
            face_filename = f"{self.video_id}_face_{self.face_count-1:03d}.jpg"
            face_path = Path("../storage/faces") / face_filename
            pil_image.save(face_path, "JPEG", quality=95)
            
            # Add to known faces
            self.known_faces.append(face_filename)
            self.known_encodings.append(face_encoding)
            new_faces.append(face_filename)
            
        return new_faces
        
    def process_video(self):
        """Process the entire video"""
        start_time = time.time()
        
        # Extract frames
        frames = self.extract_frames(self.video_path)
        
        # Process each frame
        for i, frame in enumerate(frames):
            print(f"Processing frame {i+1}/{len(frames)}")
            self.process_faces(frame, i+1)
            
        processing_time = time.time() - start_time
        print(f"Processing complete! Found {self.face_count} unique faces in {processing_time:.2f} seconds")
        
        return {
            "unique_faces_count": self.face_count,
            "faces": [f"faces/{face}" for face in self.known_faces],
            "message": f"Successfully processed video. Found {self.face_count} unique faces.",
            "processing_time_seconds": processing_time
        }

def main():
    parser = argparse.ArgumentParser(description="Process video and extract unique faces")
    parser.add_argument("video_path", help="Path to the video file")
    parser.add_argument("--video-id", help="Unique video ID for face naming")
    parser.add_argument("--fps", type=int, default=1, help="Frames per second to extract (default: 1)")
    parser.add_argument("--threshold", type=float, default=0.6, help="Face similarity threshold (default: 0.6)")
    
    args = parser.parse_args()
    
    if not os.path.exists(args.video_path):
        print(json.dumps({"error": "Video file not found"}))
        sys.exit(1)
        
    try:
        processor = FaceProcessor(args.video_path, args.video_id, args.fps, args.threshold)
        result = processor.process_video()
        
        sys.stdout.flush()  # Clear any buffered output
        print(json.dumps(result, indent=2))
        sys.stdout.flush()  # Ensure output is sent
        
    except Exception as e:
        error_response = {
            "error": f"Processing failed: {str(e)}",
            "unique_faces_count": 0,
            "faces": [],
            "message": "Video processing failed",
            "processing_time_seconds": 0
        }
        sys.stdout.flush()  # Clear any buffered output
        print(json.dumps(error_response, indent=2))
        sys.stdout.flush()  # Ensure output is sent
        sys.exit(1)

if __name__ == "__main__":
    main() 