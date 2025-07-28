#!/usr/bin/env python3
"""
Video Face Detection and Recognition System
High-accuracy face detection and deduplication using face_recognition library
"""

import sys
import json
import os
import cv2
import numpy as np
import face_recognition
from PIL import Image
import time
from pathlib import Path
from typing import List, Tuple, Dict
import argparse
import warnings

# Suppress all warnings to ensure clean JSON output
warnings.filterwarnings("ignore")

class FaceProcessor:
    def __init__(self, similarity_threshold: float = 0.6):
        """
        Initialize the face processor
        
        Args:
            similarity_threshold: Threshold for face similarity (0.6 = 60% similarity)
        """
        self.similarity_threshold = similarity_threshold
        self.known_face_encodings = []
        self.known_face_locations = []
        self.unique_faces = []
        self.face_count = 0
        
    def extract_frames(self, video_path: str, fps: int = 1) -> List[np.ndarray]:
        """
        Extract frames from video at specified FPS
        
        Args:
            video_path: Path to the video file
            fps: Frames per second to extract (default: 1)
            
        Returns:
            List of frames as numpy arrays
        """
        frames = []
        cap = cv2.VideoCapture(video_path)
        
        if not cap.isOpened():
            raise ValueError(f"Could not open video file: {video_path}")
        
        # Get video properties
        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        video_fps = cap.get(cv2.CAP_PROP_FPS)
        duration = total_frames / video_fps
        
        print(f"Video info: {total_frames} frames, {video_fps:.2f} fps, {duration:.2f}s duration")
        
        # Calculate frame interval
        frame_interval = int(video_fps / fps)
        
        frame_count = 0
        while True:
            ret, frame = cap.read()
            if not ret:
                break
                
            # Extract frame at specified interval
            if frame_count % frame_interval == 0:
                # Convert BGR to RGB
                frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
                frames.append(frame_rgb)
                
            frame_count += 1
            
        cap.release()
        print(f"Extracted {len(frames)} frames at {fps} fps")
        return frames
    
    def detect_faces_in_frame(self, frame: np.ndarray) -> Tuple[List, List]:
        """
        Detect faces in a single frame
        
        Args:
            frame: RGB frame as numpy array
            
        Returns:
            Tuple of (face_locations, face_encodings)
        """
        # Detect face locations
        face_locations = face_recognition.face_locations(frame, model="hog")
        
        # Get face encodings
        face_encodings = face_recognition.face_encodings(frame, face_locations)
        
        return face_locations, face_encodings
    
    def is_new_face(self, face_encoding: np.ndarray) -> bool:
        """
        Check if a face is new (not already seen)
        
        Args:
            face_encoding: Face encoding to check
            
        Returns:
            True if face is new, False if already seen
        """
        if len(self.known_face_encodings) == 0:
            return True
            
        # Compare with all known faces
        face_distances = face_recognition.face_distance(self.known_face_encodings, face_encoding)
        
        # If any distance is below threshold, face is not new
        return not any(face_distances <= self.similarity_threshold)
    
    def crop_and_save_face(self, frame: np.ndarray, face_location: Tuple, face_id: int) -> str:
        """
        Crop face from frame and save to disk
        
        Args:
            frame: RGB frame
            face_location: Face location (top, right, bottom, left)
            face_id: Unique face identifier
            
        Returns:
            Path to saved face image
        """
        top, right, bottom, left = face_location
        
        # Crop face with some padding
        padding = 20
        top = max(0, top - padding)
        left = max(0, left - padding)
        bottom = min(frame.shape[0], bottom + padding)
        right = min(frame.shape[1], right + padding)
        
        face_image = frame[top:bottom, left:right]
        
        # Convert to PIL Image
        pil_image = Image.fromarray(face_image)
        
        # Create faces directory if it doesn't exist
        faces_dir = Path("faces")
        faces_dir.mkdir(exist_ok=True)
        
        # Save face image
        face_filename = f"face_{face_id:03d}.jpg"
        face_path = faces_dir / face_filename
        
        pil_image.save(face_path, "JPEG", quality=95)
        
        return str(face_path)
    
    def process_video(self, video_path: str, fps: int = 1) -> Dict:
        """
        Process video and extract unique faces
        
        Args:
            video_path: Path to video file
            fps: Frames per second to process
            
        Returns:
            Dictionary with processing results
        """
        start_time = time.time()
        
        print(f"Processing video: {video_path}")
        
        # Extract frames
        frames = self.extract_frames(video_path, fps)
        
        if not frames:
            raise ValueError("No frames extracted from video")
        
        # Process each frame
        for frame_idx, frame in enumerate(frames):
            print(f"Processing frame {frame_idx + 1}/{len(frames)}")
            
            # Detect faces in frame
            face_locations, face_encodings = self.detect_faces_in_frame(frame)
            
            print(f"Found {len(face_encodings)} faces in frame {frame_idx + 1}")
            
            # Process each detected face
            for face_idx, (face_location, face_encoding) in enumerate(zip(face_locations, face_encodings)):
                
                # Check if this is a new face
                if self.is_new_face(face_encoding):
                    print(f"New face detected! Face #{self.face_count + 1}")
                    
                    # Add to known faces
                    self.known_face_encodings.append(face_encoding)
                    self.known_face_locations.append(face_location)
                    
                    # Crop and save face
                    face_path = self.crop_and_save_face(frame, face_location, self.face_count)
                    
                    # Add to unique faces list
                    self.unique_faces.append(face_path)
                    self.face_count += 1
                else:
                    print(f"Duplicate face detected (skipping)")
        
        processing_time = time.time() - start_time
        
        # Prepare response
        response = {
            "unique_faces_count": self.face_count,
            "faces": self.unique_faces,
            "message": f"Successfully processed video. Found {self.face_count} unique faces.",
            "processing_time_seconds": processing_time
        }
        
        print(f"Processing complete! Found {self.face_count} unique faces in {processing_time:.2f} seconds")
        
        return response

def main():
    """Main function to process video from command line"""
    parser = argparse.ArgumentParser(description="Process video and extract unique faces")
    parser.add_argument("video_path", help="Path to the video file")
    parser.add_argument("--fps", type=int, default=1, help="Frames per second to extract (default: 1)")
    parser.add_argument("--threshold", type=float, default=0.6, help="Face similarity threshold (default: 0.6)")
    
    args = parser.parse_args()
    
    # Validate video file exists
    if not os.path.exists(args.video_path):
        print(json.dumps({
            "error": f"Video file not found: {args.video_path}"
        }))
        sys.exit(1)
    
    try:
        # Initialize face processor
        processor = FaceProcessor(similarity_threshold=args.threshold)
        
        # Process video
        result = processor.process_video(args.video_path, args.fps)
        
        # Output JSON result (ensure clean output)
        sys.stdout.flush()  # Clear any buffered output
        print(json.dumps(result, indent=2))
        sys.stdout.flush()  # Ensure output is sent
        
    except Exception as e:
        error_response = {
            "error": f"Processing failed: {str(e)}",
            "unique_faces_count": 0,
            "faces": [],
            "message": f"Error: {str(e)}"
        }
        sys.stdout.flush()  # Clear any buffered output
        print(json.dumps(error_response, indent=2))
        sys.stdout.flush()  # Ensure output is sent
        sys.exit(1)

if __name__ == "__main__":
    main() 