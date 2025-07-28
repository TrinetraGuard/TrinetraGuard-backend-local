#!/usr/bin/env python3
"""
Comprehensive verification script for the Video Analysis System
"""

import os
import sys
import subprocess
import json

def check_go_installation():
    """Check if Go is installed and working"""
    try:
        result = subprocess.run(['go', 'version'], capture_output=True, text=True)
        if result.returncode == 0:
            print("✅ Go is installed:", result.stdout.strip())
            return True
        else:
            print("❌ Go installation check failed")
            return False
    except FileNotFoundError:
        print("❌ Go is not installed or not in PATH")
        return False

def check_python_dependencies():
    """Check Python dependencies"""
    dependencies = [
        ('cv2', 'OpenCV'),
        ('numpy', 'NumPy'),
        ('json', 'JSON'),
        ('os', 'OS'),
        ('sys', 'SYS')
    ]
    
    all_good = True
    for module, name in dependencies:
        try:
            __import__(module)
            print(f"✅ {name} - OK")
        except ImportError:
            print(f"❌ {name} - NOT FOUND")
            all_good = False
    
    return all_good

def check_file_structure():
    """Check if all required files exist"""
    required_files = [
        'backend/main.go',
        'backend/go.mod',
        'python/analyze_video_simple.py',
        'python/requirements.txt',
        'python/test_dependencies.py',
        'frontend/index.html',
        'README.md',
        'start.sh'
    ]
    
    all_good = True
    for file_path in required_files:
        if os.path.exists(file_path):
            print(f"✅ {file_path} - EXISTS")
        else:
            print(f"❌ {file_path} - MISSING")
            all_good = False
    
    return all_good

def check_directories():
    """Check if required directories exist or can be created"""
    required_dirs = [
        'backend/videos',
        'backend/faces'
    ]
    
    all_good = True
    for dir_path in required_dirs:
        if os.path.exists(dir_path):
            print(f"✅ {dir_path} - EXISTS")
        else:
            try:
                os.makedirs(dir_path, exist_ok=True)
                print(f"✅ {dir_path} - CREATED")
            except Exception as e:
                print(f"❌ {dir_path} - CANNOT CREATE: {e}")
                all_good = False
    
    return all_good

def test_go_dependencies():
    """Test Go dependencies"""
    try:
        os.chdir('backend')
        result = subprocess.run(['go', 'mod', 'tidy'], capture_output=True, text=True)
        if result.returncode == 0:
            print("✅ Go dependencies - OK")
            os.chdir('..')
            return True
        else:
            print("❌ Go dependencies - FAILED")
            print(result.stderr)
            os.chdir('..')
            return False
    except Exception as e:
        print(f"❌ Go dependencies test failed: {e}")
        os.chdir('..')
        return False

def test_python_script():
    """Test if Python script can be executed"""
    try:
        script_path = 'python/analyze_video_simple.py'
        if os.path.exists(script_path):
            # Check if script is executable
            if os.access(script_path, os.X_OK):
                print("✅ Python script - EXECUTABLE")
            else:
                print("⚠️  Python script - NOT EXECUTABLE (will be fixed)")
                os.chmod(script_path, 0o755)
                print("✅ Python script - MADE EXECUTABLE")
            return True
        else:
            print("❌ Python script - NOT FOUND")
            return False
    except Exception as e:
        print(f"❌ Python script test failed: {e}")
        return False

def main():
    """Run all verification checks"""
    print("🔍 Video Analysis System - Setup Verification")
    print("=" * 50)
    
    checks = [
        ("Go Installation", check_go_installation),
        ("File Structure", check_file_structure),
        ("Directories", check_directories),
        ("Go Dependencies", test_go_dependencies),
        ("Python Dependencies", check_python_dependencies),
        ("Python Script", test_python_script)
    ]
    
    results = []
    for check_name, check_func in checks:
        print(f"\n📋 {check_name}:")
        print("-" * 30)
        result = check_func()
        results.append((check_name, result))
    
    # Summary
    print("\n" + "=" * 50)
    print("📊 VERIFICATION SUMMARY")
    print("=" * 50)
    
    passed = sum(1 for _, result in results if result)
    total = len(results)
    
    for check_name, result in results:
        status = "✅ PASS" if result else "❌ FAIL"
        print(f"{status} - {check_name}")
    
    print(f"\nOverall: {passed}/{total} checks passed")
    
    if passed == total:
        print("\n🎉 All checks passed! The system is ready to run.")
        print("\nTo start the system:")
        print("   ./start.sh")
        print("\nOr manually:")
        print("   cd backend")
        print("   go run main.go")
        print("\nThen open: http://localhost:8080")
    else:
        print("\n❌ Some checks failed. Please fix the issues before running the system.")
        print("\nCommon fixes:")
        print("1. Install missing dependencies:")
        print("   pip install -r python/requirements.txt")
        print("2. Install system dependencies (macOS):")
        print("   brew install cmake dlib")
        print("3. Install system dependencies (Ubuntu):")
        print("   sudo apt-get install cmake libdlib-dev")

if __name__ == "__main__":
    main() 