# Flutter Integration Quick Reference

## 🚀 Quick Start

### 1. Add Dependencies
```yaml
dependencies:
  http: ^1.1.0
  image_picker: ^1.0.4
  permission_handler: ^11.0.1
```

### 2. Base URL Configuration
```dart
// Development
static const String baseUrl = 'http://localhost:8080/api/v1';

// Production
static const String baseUrl = 'https://your-domain.com/api/v1';
```

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/lost-persons/` | Create lost person report |
| `GET` | `/lost-persons/` | Get all reports |
| `GET` | `/lost-persons/:id` | Get specific report |
| `PUT` | `/lost-persons/:id` | Update report |
| `DELETE` | `/lost-persons/:id` | Delete report |
| `GET` | `/lost-persons/search?q=query` | Search reports |
| `GET` | `/images/:filename` | Get image file |

## 📝 Required Fields

### Create/Update Report
- ✅ **Name** (required)
- ✅ **Aadhar Number** (required)
- ✅ **Place Lost** (required)
- ✅ **Permanent Address** (required)
- ✅ **Image** (required)
- ⚪ **Contact Number** (optional)

## 🖼️ Image Upload

### File Requirements
- **Max Size**: 10MB
- **Formats**: JPEG, JPG, PNG, GIF
- **Field Name**: `image`

### Example Upload Code
```dart
var request = http.MultipartRequest(
  'POST',
  Uri.parse('$baseUrl/lost-persons/'),
);

request.fields['name'] = name;
request.fields['aadhar_number'] = aadharNumber;
request.fields['place_lost'] = placeLost;
request.fields['permanent_address'] = permanentAddress;

request.files.add(
  await http.MultipartFile.fromPath('image', imageFile.path),
);
```

## 📊 Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "aadhar_number": "123456789012",
    "contact_number": "9876543210",
    "place_lost": "Mumbai Central",
    "permanent_address": "123 Main St",
    "image_path": "uploads/filename.jpg",
    "upload_timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## 🔧 Essential Code Snippets

### 1. Create Report
```dart
var result = await LostPersonService.createLostPersonReport(
  name: "John Doe",
  aadharNumber: "123456789012",
  contactNumber: "9876543210",
  placeLost: "Mumbai Central",
  permanentAddress: "123 Main St",
  imageFile: imageFile,
);

if (result['success']) {
  // Handle success
} else {
  // Handle error: result['error']
}
```

### 2. Get All Reports
```dart
var result = await LostPersonService.getAllLostPersons();
if (result['success']) {
  var lostPersonList = LostPersonList.fromJson(result['data']);
  // Use lostPersonList.lostPersons
}
```

### 3. Search Reports
```dart
var result = await LostPersonService.searchLostPersons("John Doe");
if (result['success']) {
  var lostPersonList = LostPersonList.fromJson(result['data']);
  // Use lostPersonList.lostPersons
}
```

### 4. Display Image
```dart
String imageUrl = LostPersonService.getImageUrl(person.imagePath);
Image.network(imageUrl)
```

## ⚠️ Important Notes

### 1. Error Handling
- Always check `result['success']` before using data
- Display `result['error']` for user feedback
- Handle network timeouts

### 2. Permissions
```dart
// Request camera permission
var status = await Permission.camera.request();
if (status != PermissionStatus.granted) {
  // Handle permission denied
}
```

### 3. Image Compression
```dart
final XFile? image = await _picker.pickImage(
  source: ImageSource.camera,
  maxWidth: 1024,
  maxHeight: 1024,
  imageQuality: 85,
);
```

## 🐛 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Network error | Check if backend is running |
| Image upload fails | Verify file size < 10MB |
| Permission denied | Request camera/storage permissions |
| Image not loading | Check URL construction |
| CORS error | Backend should allow Flutter requests |

## 📱 Testing Checklist

- [ ] Test image upload with different formats
- [ ] Test with large images (>5MB)
- [ ] Test network connectivity scenarios
- [ ] Test form validation
- [ ] Test on both Android and iOS
- [ ] Test camera and gallery access
- [ ] Test search functionality
- [ ] Test error scenarios

## 🔗 Useful Links

- [Complete Integration Guide](FLUTTER_INTEGRATION_GUIDE.md)
- [API Documentation](API_DOCUMENTATION.md)
- [Backend Repository](README.md)

---

**Need Help?** Check the complete integration guide for detailed examples and troubleshooting.
