#!/usr/bin/env python3
import cv2
import numpy as np
import os

def test_face_detection():
    """Test face detection on a sample frame"""
    
    # Load the cascade classifier
    cascade_paths = [
        '/opt/homebrew/Cellar/opencv/4.11.0_1/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
        '/opt/homebrew/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
        '/usr/local/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
        '/usr/local/share/opencv4/haarcascades/haarcascade_frontalface_default.xml'
    ]
    
    face_cascade = None
    for path in cascade_paths:
        if os.path.exists(path):
            face_cascade = cv2.CascadeClassifier(path)
            print(f"✅ Loaded cascade from: {path}")
            break
    
    if face_cascade is None:
        print("❌ Could not load face cascade classifier")
        return
    
    # Test with different parameters
    test_params = [
        {"scaleFactor": 1.05, "minNeighbors": 15, "minSize": (150, 150), "maxSize": (350, 350)},
        {"scaleFactor": 1.1, "minNeighbors": 10, "minSize": (100, 100), "maxSize": (400, 400)},
        {"scaleFactor": 1.05, "minNeighbors": 8, "minSize": (80, 80), "maxSize": (500, 500)},
    ]
    
    print("\n🔍 Testing face detection parameters:")
    for i, params in enumerate(test_params):
        print(f"\nTest {i+1}: {params}")
        
        # Test with a sample video frame
        cap = cv2.VideoCapture("../backend/videos/1753685901_VID20250220152617.mp4")
        if not cap.isOpened():
            print("❌ Could not open video file")
            return
        
        # Read a frame from the middle of the video
        cap.set(cv2.CAP_PROP_POS_FRAMES, 100)
        ret, frame = cap.read()
        cap.release()
        
        if not ret:
            print("❌ Could not read frame")
            continue
        
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        
        faces = face_cascade.detectMultiScale(
            gray,
            scaleFactor=params["scaleFactor"],
            minNeighbors=params["minNeighbors"],
            minSize=params["minSize"],
            maxSize=params["maxSize"]
        )
        
        print(f"   Found {len(faces)} faces")
        
        # Draw rectangles around detected faces
        for (x, y, w, h) in faces:
            cv2.rectangle(frame, (x, y), (x+w, y+h), (255, 0, 0), 2)
        
        # Save the frame with detected faces
        output_path = f"test_frame_{i+1}.jpg"
        cv2.imwrite(output_path, frame)
        print(f"   Saved frame to: {output_path}")

if __name__ == "__main__":
    test_face_detection() 