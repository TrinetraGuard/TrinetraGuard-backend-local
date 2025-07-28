#!/usr/bin/env python3
import cv2
import numpy as np
import json
import os
import sys
from datetime import timedelta
import math
import time

def log_progress(message, flush=True):
    """Log progress message with timestamp"""
    timestamp = time.strftime("%H:%M:%S")
    print(f"[{timestamp}] {message}", flush=flush)

def calculate_face_similarity(face1, face2):
    """Calculate similarity between two face images using multiple methods - improved for accuracy"""
    try:
        # Resize both faces to same size for comparison
        size = (150, 150)  # Larger size for better comparison
        face1_resized = cv2.resize(face1, size)
        face2_resized = cv2.resize(face2, size)
        
        # Convert to grayscale
        gray1 = cv2.cvtColor(face1_resized, cv2.COLOR_BGR2GRAY)
        gray2 = cv2.cvtColor(face2_resized, cv2.COLOR_BGR2GRAY)
        
        # Method 1: Template matching with multiple scales
        max_similarity = 0
        for scale in [0.8, 0.9, 1.0, 1.1, 1.2]:
            scaled_gray2 = cv2.resize(gray2, None, fx=scale, fy=scale)
            if scaled_gray2.shape[0] <= gray1.shape[0] and scaled_gray2.shape[1] <= gray1.shape[1]:
                result = cv2.matchTemplate(gray1, scaled_gray2, cv2.TM_CCOEFF_NORMED)
                max_val = np.max(result)
                max_similarity = max(max_similarity, max_val)
        
        # Method 2: Histogram comparison
        hist1 = cv2.calcHist([gray1], [0], None, [256], [0, 256])
        hist2 = cv2.calcHist([gray2], [0], None, [256], [0, 256])
        hist_similarity = cv2.compareHist(hist1, hist2, cv2.HISTCMP_CORREL)
        
        # Method 3: Structural similarity (SSIM-like)
        diff = cv2.absdiff(gray1, gray2)
        ssim_similarity = 1.0 - (np.mean(diff) / 255.0)
        
        # Method 4: Edge comparison
        edges1 = cv2.Canny(gray1, 50, 150)
        edges2 = cv2.Canny(gray2, 50, 150)
        edge_similarity = 1.0 - (np.sum(cv2.absdiff(edges1, edges2)) / (edges1.shape[0] * edges1.shape[1] * 255))
        
        # Method 5: Feature-based comparison using ORB
        orb = cv2.ORB_create()
        kp1, des1 = orb.detectAndCompute(gray1, None)
        kp2, des2 = orb.detectAndCompute(gray2, None)
        
        if des1 is not None and des2 is not None and len(des1) > 0 and len(des2) > 0:
            bf = cv2.BFMatcher(cv2.NORM_HAMMING, crossCheck=True)
            matches = bf.match(des1, des2)
            if len(matches) > 0:
                feature_similarity = len(matches) / max(len(kp1), len(kp2))
            else:
                feature_similarity = 0
        else:
            feature_similarity = 0
        
        # Combine all similarities with improved weighting
        combined_similarity = (
            max_similarity * 0.35 +      # Template matching (highest weight)
            hist_similarity * 0.25 +     # Histogram comparison
            ssim_similarity * 0.20 +     # Structural similarity
            edge_similarity * 0.15 +     # Edge comparison
            feature_similarity * 0.05    # Feature-based comparison
        )
        
        return combined_similarity
    except:
        return 0.0

def calculate_face_quality(face_img):
    """Calculate quality score for a face image - extremely strict for 100% accuracy"""
    try:
        height, width = face_img.shape[:2]
        
        # Convert to grayscale for analysis
        gray = cv2.cvtColor(face_img, cv2.COLOR_BGR2GRAY)
        
        # 1. Sharpness (Laplacian variance) - very strict
        laplacian_var = cv2.Laplacian(gray, cv2.CV_64F).var()
        sharpness_score = min(laplacian_var / 100, 100)  # Normalize to 0-100
        
        # 2. Contrast - very strict
        contrast = np.std(gray)
        contrast_score = min(contrast / 2, 100)
        
        # 3. Brightness - very strict
        brightness = np.mean(gray)
        brightness_score = 100 - abs(brightness - 128) / 1.28  # Penalize too bright/dark
        
        # 4. Size score - very strict
        size_score = min((width * height) / 50000, 100)  # Prefer larger faces
        
        # 5. Aspect ratio - very strict
        aspect_ratio = width / height
        if 0.8 <= aspect_ratio <= 1.2:  # Very strict aspect ratio
            aspect_score = 100
        else:
            aspect_score = 0
        
        # 6. Face completeness check - very strict
        completeness_score = 0
        if is_complete_human_face(face_img):
            completeness_score = 100
        
        # 7. Edge density - very strict
        edges = cv2.Canny(gray, 50, 150)
        edge_density = np.sum(edges > 0) / (width * height)
        if edge_density < 0.1:  # Very low edge density for clean faces
            edge_score = 100
        else:
            edge_score = 0
        
        # Combine all scores with strict weighting
        total_score = (
            sharpness_score * 0.25 +
            contrast_score * 0.20 +
            brightness_score * 0.15 +
            size_score * 0.15 +
            aspect_score * 0.10 +
            completeness_score * 0.10 +
            edge_score * 0.05
        )
        
        return total_score
    except:
        return 0

def is_complete_human_face(face_img):
    """Check if the detected face is a complete human face - very lenient for different video types"""
    try:
        height, width = face_img.shape[:2]
        
        # Check aspect ratio (human faces are roughly 1:1 to 1:1.3) - very lenient
        aspect_ratio = width / height
        if aspect_ratio < 0.7 or aspect_ratio > 1.5: # Very lenient aspect ratio
            return False
        
        # Check size (must be large enough to be a complete face) - very lenient
        if width < 120 or height < 120: # Very lenient minimum size
            return False
        
        # Check if face is too large (likely a false positive) - very lenient
        if width > 450 or height > 450: # Very lenient maximum size
            return False
        
        # Convert to HSV and check skin tone - very lenient
        hsv = cv2.cvtColor(face_img, cv2.COLOR_BGR2HSV)
        
        # Define skin tone ranges (very lenient)
        lower_skin = np.array([0, 30, 80], dtype=np.uint8) # Very lenient skin tone
        upper_skin = np.array([25, 255, 255], dtype=np.uint8)
        
        # Create mask for skin tones
        skin_mask = cv2.inRange(hsv, lower_skin, upper_skin)
        skin_pixels = cv2.countNonZero(skin_mask)
        total_pixels = width * height
        skin_ratio = skin_pixels / total_pixels
        
        # Must have reasonable skin tone ratio (not too much, not too little) - very lenient
        if skin_ratio < 0.15 or skin_ratio > 0.8: # Very lenient skin ratio
            return False
        
        # Check edge density (faces should have moderate edges) - very lenient
        gray = cv2.cvtColor(face_img, cv2.COLOR_BGR2GRAY)
        edges = cv2.Canny(gray, 50, 150)
        edge_density = np.sum(edges > 0) / (height * width)
        if edge_density > 0.25: # Very lenient edge density
            return False
        
        # Check symmetry (human faces are roughly symmetrical) - very lenient
        left_half = gray[:, :width//2]
        right_half = cv2.flip(gray[:, width//2:], 1)
        if right_half.shape[1] != left_half.shape[1]:
            right_half = right_half[:, :left_half.shape[1]]
        symmetry = np.corrcoef(left_half.flatten(), right_half.flatten())[0, 1]
        if symmetry < 0.3: # Very lenient symmetry
            return False
        
        # Additional check: look for eyes using eye cascade - very lenient
        eye_cascade = cv2.CascadeClassifier('/opt/homebrew/Cellar/opencv/4.11.0_1/share/opencv4/haarcascades/haarcascade_eye.xml')
        eyes = eye_cascade.detectMultiScale(gray, 1.1, 2)  # More lenient eye detection
        
        # Must have at least 1 eye detected to be a complete face - very lenient
        if len(eyes) < 1:
            return False
        
        return True
    except Exception as e:
        return False

def analyze_video(video_path, faces_dir):
    """Analyze video and detect only complete, unique human faces"""
    
    log_progress("🚀 Starting video analysis...")
    
    # Create faces directory if it doesn't exist
    os.makedirs(faces_dir, exist_ok=True)
    log_progress("📁 Created faces directory")
    
    # Load the face cascade classifier
    try:
        # Try different possible paths for the cascade file
        cascade_paths = [
            '/opt/homebrew/Cellar/opencv/4.11.0_1/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
            '/opt/homebrew/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
            '/usr/local/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
            '/usr/share/opencv4/haarcascades/haarcascade_frontalface_default.xml',
            'haarcascade_frontalface_default.xml'
        ]
        
        face_cascade = None
        for path in cascade_paths:
            try:
                face_cascade = cv2.CascadeClassifier(path)
                if not face_cascade.empty():
                    log_progress(f"🔍 Loaded face detection model from: {path}")
                    break
            except:
                continue
        
        if face_cascade is None or face_cascade.empty():
            log_progress("❌ Error: Could not load face detection model")
            return None
            
    except Exception as e:
        log_progress(f"❌ Error loading face detection model: {str(e)}")
        return None
    
    # Open the video
    cap = cv2.VideoCapture(video_path)
    if not cap.isOpened():
        log_progress("❌ Error: Could not open video")
        return None
    
    # Get video properties
    fps = cap.get(cv2.CAP_PROP_FPS)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    duration = total_frames / fps if fps > 0 else 0
    
    log_progress(f"📹 Video info: {total_frames} frames, {fps:.1f} FPS, {duration:.1f}s duration")
    
    # Store detected faces
    detected_faces = []
    frame_count = 0
    processed_frames = 0
    
    log_progress("🔍 Starting face detection (processing every 5th frame to catch all faces)...")
    
    while True:
        ret, frame = cap.read()
        if not ret:
            break
            
        frame_count += 1
        
        # Process every 5th frame to catch all faces
        if frame_count % 5 != 0:
            continue
        
        processed_frames += 1
        
        # Convert to grayscale for face detection
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        
        # Detect faces with very sensitive parameters to catch faces in different video types
        faces = face_cascade.detectMultiScale(
            gray,
            scaleFactor=1.02,     # Very sensitive scaling to catch different face types
            minNeighbors=5,       # Very low threshold to detect more faces
            minSize=(100, 100),   # Very small minimum to catch various face sizes
            maxSize=(500, 500)    # Very large maximum to catch different faces
        )
        
        # Debug: Log if no faces found in this frame
        if len(faces) == 0 and frame_count % 50 == 0:  # Log every 50th frame
            log_progress(f"Debug: No faces detected in frame {frame_count}", flush=True)
        
        log_progress(f"Frame {frame_count}: Found {len(faces)} raw faces", flush=True)
        
        # Process each detected face
        for (x, y, w, h) in faces:
            # Extract face region
            face_img = frame[y:y+h, x:x+w]
            
            # Check if it's a complete human face
            if not is_complete_human_face(face_img):
                continue
            
            # Calculate quality score
            quality_score = calculate_face_quality(face_img)
            
            # Only keep reasonable-quality faces to catch faces in different video types
            if quality_score < 15:  # Very low threshold to catch more faces
                log_progress(f"Face rejected: quality score {quality_score:.1f} < 15", flush=True)
                continue
            
            # Calculate timestamp
            timestamp_seconds = frame_count / fps
            timestamp = str(timedelta(seconds=int(timestamp_seconds)))
            
            # Check if this face is similar to any existing face and keep highest quality
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
                        
                        # If similarity is high, it's a duplicate - much less strict to catch all unique people
                        if similarity > 0.75:  # Much higher threshold to avoid rejecting different people
                            log_progress(f"Face rejected: similarity {similarity:.3f} > 0.75 (duplicate)", flush=True)
                            is_duplicate = True
                            break
            
            if is_duplicate:
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
                
                log_progress(f"👤 Found new person #{face_index} at {timestamp}")
        
        # Progress update every 10 processed frames
        if processed_frames % 10 == 0:
            progress = (frame_count / total_frames) * 100
            log_progress(f"📊 Progress: {progress:.1f}% - Found {len(detected_faces)} unique faces so far")
    
    # Release video capture
    cap.release()
    
    log_progress(f"✅ Analysis complete! Found {len(detected_faces)} unique faces")
    
    # Sort faces by quality and limit to top results
    detected_faces.sort(key=lambda x: x.get('quality', 0), reverse=True)
    
    # Limit to maximum 25 unique faces to catch all people in video
    max_faces = min(len(detected_faces), 25)
    detected_faces = detected_faces[:max_faces]
    
    log_progress(f"🎯 Final result: {max_faces} high-quality unique faces selected")
    
    # Output final results as clean JSON on a single line
    result = {
        "total_people": len(detected_faces),
        "faces": []
    }
    
    for i, face_data in enumerate(detected_faces):
        result["faces"].append({
            "image": face_data['image'],
            "timestamps": face_data['timestamps']
        })
    
    # Output clean JSON on a single line to avoid parsing issues
    print(json.dumps(result, separators=(',', ':')))
    sys.stdout.flush()  # Ensure output is flushed immediately
    
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