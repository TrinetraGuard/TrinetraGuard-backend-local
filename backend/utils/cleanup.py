#!/usr/bin/env python3
"""
Cleanup utility for video processing backend
Removes old video files and face images to manage disk space
"""

import os
import shutil
import time
from pathlib import Path
import argparse

def cleanup_old_files(directory: str, max_age_hours: int = 24):
    """
    Remove files older than specified hours
    
    Args:
        directory: Directory to clean
        max_age_hours: Maximum age in hours before deletion
    """
    if not os.path.exists(directory):
        print(f"Directory {directory} does not exist")
        return
    
    current_time = time.time()
    max_age_seconds = max_age_hours * 3600
    
    removed_count = 0
    total_size = 0
    
    for filename in os.listdir(directory):
        file_path = os.path.join(directory, filename)
        
        if os.path.isfile(file_path):
            file_age = current_time - os.path.getmtime(file_path)
            
            if file_age > max_age_seconds:
                file_size = os.path.getsize(file_path)
                total_size += file_size
                
                try:
                    os.remove(file_path)
                    removed_count += 1
                    print(f"Removed: {file_path}")
                except OSError as e:
                    print(f"Error removing {file_path}: {e}")
    
    if removed_count > 0:
        print(f"Cleaned {directory}: {removed_count} files removed, {total_size / (1024*1024):.2f} MB freed")
    else:
        print(f"No files to clean in {directory}")

def main():
    parser = argparse.ArgumentParser(description="Clean up old video and face files")
    parser.add_argument("--videos", action="store_true", help="Clean videos directory")
    parser.add_argument("--faces", action="store_true", help="Clean faces directory")
    parser.add_argument("--all", action="store_true", help="Clean both directories")
    parser.add_argument("--age", type=int, default=24, help="Maximum age in hours (default: 24)")
    
    args = parser.parse_args()
    
    if not any([args.videos, args.faces, args.all]):
        print("Please specify --videos, --faces, or --all")
        return
    
    if args.all or args.videos:
        cleanup_old_files("videos", args.age)
    
    if args.all or args.faces:
        cleanup_old_files("faces", args.age)

if __name__ == "__main__":
    main() 