#!/usr/bin/env python3
import cv2
import numpy as np
import os

def create_test_video():
    """Create a test video with multiple faces for testing"""
    
    # Video parameters
    width, height = 640, 480
    fps = 30
    duration = 5  # 5 seconds
    total_frames = fps * duration
    
    # Create video writer
    fourcc = cv2.VideoWriter_fourcc(*'mp4v')
    out = cv2.VideoWriter('test_video.mp4', fourcc, fps, (width, height))
    
    # Create different colored rectangles to simulate faces
    face_positions = [
        (100, 100, 150, 150, (255, 0, 0)),    # Blue face
        (300, 100, 150, 150, (0, 255, 0)),    # Green face
        (500, 100, 150, 150, (0, 0, 255)),    # Red face
        (200, 300, 150, 150, (255, 255, 0)),  # Yellow face
        (400, 300, 150, 150, (255, 0, 255)),  # Magenta face
    ]
    
    for frame_num in range(total_frames):
        # Create black background
        frame = np.zeros((height, width, 3), dtype=np.uint8)
        
        # Add some faces that appear at different times
        for i, (x, y, w, h, color) in enumerate(face_positions):
            # Make faces appear at different times
            if frame_num > (i * fps):  # Each face appears after i seconds
                cv2.rectangle(frame, (x, y), (x + w, y + h), color, -1)
                cv2.putText(frame, f'Person {i+1}', (x, y-10), 
                           cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
        
        out.write(frame)
    
    out.release()
    print("Test video created: test_video.mp4")
    print("This video contains 5 simulated faces that appear at different times")

if __name__ == "__main__":
    create_test_video() 