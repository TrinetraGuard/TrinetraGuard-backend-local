#!/usr/bin/env python3
"""
Test script to verify all Python dependencies are properly installed
"""

def test_dependencies():
    """Test if all required packages are available"""
    try:
        import cv2
        print("✅ OpenCV (cv2) - OK")
    except ImportError:
        print("❌ OpenCV (cv2) - NOT FOUND")
        return False
    
    try:
        import face_recognition
        print("✅ face_recognition - OK")
    except ImportError:
        print("❌ face_recognition - NOT FOUND")
        print("   Install with: pip install face-recognition")
        return False
    
    try:
        import numpy as np
        print("✅ NumPy - OK")
    except ImportError:
        print("❌ NumPy - NOT FOUND")
        return False
    
    try:
        import json
        print("✅ JSON - OK")
    except ImportError:
        print("❌ JSON - NOT FOUND")
        return False
    
    try:
        import os
        import sys
        print("✅ OS/SYS - OK")
    except ImportError:
        print("❌ OS/SYS - NOT FOUND")
        return False
    
    print("\n🎉 All dependencies are properly installed!")
    return True

if __name__ == "__main__":
    print("Testing Python dependencies...\n")
    success = test_dependencies()
    
    if success:
        print("\n✅ Ready to run video analysis!")
    else:
        print("\n❌ Please install missing dependencies before running the system.")
        print("Run: pip install -r requirements.txt") 