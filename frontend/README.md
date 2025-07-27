# Video Analysis Service - Frontend

A modern, responsive web interface for the Video Analysis Service API. This frontend provides an intuitive way to interact with all the service features including video upload, analysis, and person finding.

## 🚀 Features

- **📹 Video Management**: Upload, view, and manage video files
- **🔍 Analysis Operations**: Start and monitor video analysis jobs
- **👤 Person Finder**: Upload reference images and search for people across videos
- **💚 Health Monitoring**: Check service status and health
- **📱 Responsive Design**: Works on desktop, tablet, and mobile devices
- **🎨 Modern UI**: Clean, intuitive interface with smooth animations
- **🌙 Dark Mode Support**: Automatic dark mode detection
- **⚡ Real-time Updates**: Live status updates and progress tracking

## 📁 File Structure

```
frontend/
├── index.html          # Main HTML file
├── styles.css          # Additional CSS styles
└── README.md           # This file
```

## 🛠️ Setup Instructions

### Prerequisites
- Video Analysis Service backend running on `http://localhost:8080`
- Modern web browser (Chrome, Firefox, Safari, Edge)

### Quick Start

1. **Start the Backend Service**
   ```bash
   # Navigate to the backend directory
   cd /path/to/video-analysis-service
   
   # Build and run the service
   make build
   make run
   ```

2. **Open the Frontend**
   ```bash
   # Navigate to the frontend directory
   cd frontend
   
   # Open in your browser (or use a local server)
   open index.html
   ```

   **Alternative**: Use a local server for better development experience
   ```bash
   # Using Python 3
   python -m http.server 3000
   
   # Using Node.js (if you have http-server installed)
   npx http-server -p 3000
   
   # Then open http://localhost:3000 in your browser
   ```

## 🎯 Usage Guide

### 1. Video Management Tab

#### Upload a Video
1. Click on the **📹 Videos** tab
2. Click **Choose video file** or drag and drop a video file
3. Click **Upload Video**
4. Wait for the upload to complete
5. The video will appear in the **Uploaded Videos** section

#### View Video Details
- Each video card shows:
  - Original filename
  - File size
  - Upload status
  - Duration (if available)
  - Resolution (if available)
  - Upload date

#### Download or Delete Videos
- Use the **Download** button to download the original video file
- Use the **Delete** button to remove the video and its analysis data

### 2. Analysis Tab

#### Start Video Analysis
1. Click on the **🔍 Analysis** tab
2. Select a video from the dropdown
3. Click **Start Analysis**
4. Monitor progress with **Check Status**
5. View results with **Get Results**

#### Analysis Results Include
- Total frames analyzed
- Total people detected
- Number of unique people
- Analysis completion date

### 3. Person Finder Tab

#### Upload Reference Image
1. Click on the **👤 Person Finder** tab
2. Click **Choose reference image** or drag and drop an image
3. Add a description (optional)
4. Click **Upload Image**
5. The image will appear in the **Reference Images** section

#### Search for a Person
1. Select a reference image from the dropdown
2. Select one or more videos to search in
3. Click **Search for Person**
4. Monitor search progress with **Check Search Status**
5. View results with **Get Search Results**

#### Search Results Include
- Video ID where person was found
- First and last appearance timestamps
- Total number of appearances
- Confidence score

### 4. Health Tab

#### Check Service Status
1. Click on the **💚 Health** tab
2. Click **Check Health**
3. View service status, version, and timestamp

## 🔧 Configuration

### API Endpoint Configuration
The frontend is configured to connect to the backend at `http://localhost:8080`. To change this:

1. Open `index.html`
2. Find the line: `const API_BASE = 'http://localhost:8080/api/v1';`
3. Change the URL to match your backend configuration

### CORS Configuration
If you're running the frontend on a different port, ensure CORS is properly configured in the backend.

## 🎨 Customization

### Styling
- Modify `styles.css` for custom styling
- The interface uses CSS Grid and Flexbox for responsive design
- Color scheme can be customized by changing CSS variables

### Functionality
- All API interactions are handled in the JavaScript section of `index.html`
- Functions are well-commented and organized by feature
- Easy to extend with additional features

## 📱 Browser Support

- **Chrome**: 60+
- **Firefox**: 55+
- **Safari**: 12+
- **Edge**: 79+

## 🐛 Troubleshooting

### Common Issues

#### "Failed to fetch" Error
- Ensure the backend service is running
- Check that the API endpoint URL is correct
- Verify CORS configuration

#### File Upload Issues
- Check file size limits (100MB max)
- Ensure file format is supported
- Verify network connection

#### Analysis/Search Not Starting
- Check that the selected video/image exists
- Ensure the backend service is healthy
- Verify database connectivity

### Debug Mode
Open browser developer tools (F12) to see:
- Network requests and responses
- JavaScript errors
- Console logs for debugging

## 🔒 Security Notes

- The frontend runs entirely in the browser
- No sensitive data is stored locally
- All API calls include proper error handling
- File uploads are validated before sending

## 🚀 Performance Tips

- Use supported video formats for better performance
- Keep reference images under 10MB for faster processing
- Close unused browser tabs to free up memory
- Use a modern browser for best performance

## 📞 Support

For issues or questions:
1. Check the browser console for error messages
2. Verify the backend service is running and healthy
3. Review the API documentation
4. Check the backend logs for server-side issues

## 🔄 Updates

To update the frontend:
1. Replace the `index.html` file
2. Update `styles.css` if needed
3. Clear browser cache if changes don't appear
4. Test all functionality after updates

---

*Frontend Version: 1.0.0*  
*Last Updated: 2025-07-27* 