# Flutter Integration Guide for TrinetraGuard Backend

## Overview

This guide provides step-by-step instructions for integrating the TrinetraGuard Lost Person Backend with your Flutter application. The backend provides RESTful API endpoints for managing lost person reports with image upload capabilities.

## Backend Information

- **Base URL**: `http://localhost:8080/api/v1` (change to your server URL)
- **Content-Type**: `multipart/form-data` for file uploads
- **Response Format**: JSON
- **Authentication**: None required (for now)

## Prerequisites

1. **Flutter SDK**: Ensure you have Flutter installed and configured
2. **HTTP Package**: Add the `http` package to your `pubspec.yaml`
3. **Image Picker**: Add the `image_picker` package for camera/gallery access
4. **File Picker**: Add the `file_picker` package for file selection

## Setup Dependencies

Add these dependencies to your `pubspec.yaml`:

```yaml
dependencies:
  flutter:
    sdk: flutter
  http: ^1.1.0
  image_picker: ^1.0.4
  file_picker: ^6.1.1
  path: ^1.8.3
  permission_handler: ^11.0.1
```

Run:
```bash
flutter pub get
```

## API Service Class

Create a service class to handle all API calls:

```dart
import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;

class LostPersonService {
  static const String baseUrl = 'http://localhost:8080/api/v1';
  
  // Create a new lost person report
  static Future<Map<String, dynamic>> createLostPersonReport({
    required String name,
    required String aadharNumber,
    String? contactNumber,
    required String placeLost,
    required String permanentAddress,
    required File imageFile,
  }) async {
    try {
      var request = http.MultipartRequest(
        'POST',
        Uri.parse('$baseUrl/lost-persons/'),
      );

      // Add form fields
      request.fields['name'] = name;
      request.fields['aadhar_number'] = aadharNumber;
      if (contactNumber != null) {
        request.fields['contact_number'] = contactNumber;
      }
      request.fields['place_lost'] = placeLost;
      request.fields['permanent_address'] = permanentAddress;

      // Add image file
      request.files.add(
        await http.MultipartFile.fromPath(
          'image',
          imageFile.path,
        ),
      );

      var response = await request.send();
      var responseData = await response.stream.bytesToString();
      var jsonResponse = json.decode(responseData);

      if (response.statusCode == 201) {
        return {
          'success': true,
          'data': jsonResponse['data'],
          'message': jsonResponse['message'],
        };
      } else {
        return {
          'success': false,
          'error': jsonResponse['error'] ?? 'Failed to create report',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': 'Network error: $e',
      };
    }
  }

  // Get all lost person reports
  static Future<Map<String, dynamic>> getAllLostPersons() async {
    try {
      var response = await http.get(
        Uri.parse('$baseUrl/lost-persons/'),
      );

      var jsonResponse = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': jsonResponse['data'],
          'message': jsonResponse['message'],
        };
      } else {
        return {
          'success': false,
          'error': jsonResponse['error'] ?? 'Failed to fetch reports',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': 'Network error: $e',
      };
    }
  }

  // Get a specific lost person report by ID
  static Future<Map<String, dynamic>> getLostPersonById(String id) async {
    try {
      var response = await http.get(
        Uri.parse('$baseUrl/lost-persons/$id'),
      );

      var jsonResponse = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': jsonResponse['data'],
          'message': jsonResponse['message'],
        };
      } else {
        return {
          'success': false,
          'error': jsonResponse['error'] ?? 'Failed to fetch report',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': 'Network error: $e',
      };
    }
  }

  // Search lost persons
  static Future<Map<String, dynamic>> searchLostPersons(String query) async {
    try {
      var response = await http.get(
        Uri.parse('$baseUrl/lost-persons/search?q=${Uri.encodeComponent(query)}'),
      );

      var jsonResponse = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': jsonResponse['data'],
          'message': jsonResponse['message'],
        };
      } else {
        return {
          'success': false,
          'error': jsonResponse['error'] ?? 'Failed to search reports',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': 'Network error: $e',
      };
    }
  }

  // Update a lost person report
  static Future<Map<String, dynamic>> updateLostPersonReport({
    required String id,
    required String name,
    required String aadharNumber,
    String? contactNumber,
    required String placeLost,
    required String permanentAddress,
    File? imageFile,
  }) async {
    try {
      var request = http.MultipartRequest(
        'PUT',
        Uri.parse('$baseUrl/lost-persons/$id'),
      );

      // Add form fields
      request.fields['name'] = name;
      request.fields['aadhar_number'] = aadharNumber;
      if (contactNumber != null) {
        request.fields['contact_number'] = contactNumber;
      }
      request.fields['place_lost'] = placeLost;
      request.fields['permanent_address'] = permanentAddress;

      // Add image file if provided
      if (imageFile != null) {
        request.files.add(
          await http.MultipartFile.fromPath(
            'image',
            imageFile.path,
          ),
        );
      }

      var response = await request.send();
      var responseData = await response.stream.bytesToString();
      var jsonResponse = json.decode(responseData);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': jsonResponse['data'],
          'message': jsonResponse['message'],
        };
      } else {
        return {
          'success': false,
          'error': jsonResponse['error'] ?? 'Failed to update report',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': 'Network error: $e',
      };
    }
  }

  // Delete a lost person report
  static Future<Map<String, dynamic>> deleteLostPersonReport(String id) async {
    try {
      var response = await http.delete(
        Uri.parse('$baseUrl/lost-persons/$id'),
      );

      var jsonResponse = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'message': jsonResponse['message'],
        };
      } else {
        return {
          'success': false,
          'error': jsonResponse['error'] ?? 'Failed to delete report',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': 'Network error: $e',
      };
    }
  }

  // Get image URL
  static String getImageUrl(String imagePath) {
    // Extract filename from path (e.g., "uploads/filename.jpg" -> "filename.jpg")
    String filename = imagePath.split('/').last;
    return '$baseUrl/images/$filename';
  }
}
```

## Data Models

Create data models for type safety:

```dart
class LostPerson {
  final String id;
  final String name;
  final String aadharNumber;
  final String? contactNumber;
  final String placeLost;
  final String permanentAddress;
  final String imagePath;
  final DateTime uploadTimestamp;

  LostPerson({
    required this.id,
    required this.name,
    required this.aadharNumber,
    this.contactNumber,
    required this.placeLost,
    required this.permanentAddress,
    required this.imagePath,
    required this.uploadTimestamp,
  });

  factory LostPerson.fromJson(Map<String, dynamic> json) {
    return LostPerson(
      id: json['id'],
      name: json['name'],
      aadharNumber: json['aadhar_number'],
      contactNumber: json['contact_number'],
      placeLost: json['place_lost'],
      permanentAddress: json['permanent_address'],
      imagePath: json['image_path'],
      uploadTimestamp: DateTime.parse(json['upload_timestamp']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'aadhar_number': aadharNumber,
      'contact_number': contactNumber,
      'place_lost': placeLost,
      'permanent_address': permanentAddress,
      'image_path': imagePath,
      'upload_timestamp': uploadTimestamp.toIso8601String(),
    };
  }
}

class LostPersonList {
  final int total;
  final List<LostPerson> lostPersons;

  LostPersonList({
    required this.total,
    required this.lostPersons,
  });

  factory LostPersonList.fromJson(Map<String, dynamic> json) {
    return LostPersonList(
      total: json['total'],
      lostPersons: (json['lost_persons'] as List)
          .map((person) => LostPerson.fromJson(person))
          .toList(),
    );
  }
}
```

## Complete Flutter App Example

Here's a complete example of how to use the API in a Flutter app:

```dart
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:permission_handler/permission_handler.dart';

class LostPersonApp extends StatefulWidget {
  @override
  _LostPersonAppState createState() => _LostPersonAppState();
}

class _LostPersonAppState extends State<LostPersonApp> {
  List<LostPerson> lostPersons = [];
  bool isLoading = false;
  final ImagePicker _picker = ImagePicker();

  @override
  void initState() {
    super.initState();
    _loadLostPersons();
  }

  Future<void> _loadLostPersons() async {
    setState(() {
      isLoading = true;
    });

    var result = await LostPersonService.getAllLostPersons();
    
    setState(() {
      isLoading = false;
    });

    if (result['success']) {
      var lostPersonList = LostPersonList.fromJson(result['data']);
      setState(() {
        lostPersons = lostPersonList.lostPersons;
      });
    } else {
      _showErrorSnackBar(result['error']);
    }
  }

  Future<void> _addLostPerson() async {
    // Request camera permission
    var status = await Permission.camera.request();
    if (status != PermissionStatus.granted) {
      _showErrorSnackBar('Camera permission is required');
      return;
    }

    // Pick image
    final XFile? image = await _picker.pickImage(
      source: ImageSource.camera,
      maxWidth: 1024,
      maxHeight: 1024,
      imageQuality: 85,
    );

    if (image == null) return;

    // Show form dialog
    _showAddPersonDialog(File(image.path));
  }

  void _showAddPersonDialog(File imageFile) {
    final nameController = TextEditingController();
    final aadharController = TextEditingController();
    final contactController = TextEditingController();
    final placeController = TextEditingController();
    final addressController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Add Lost Person Report'),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Image.file(
                imageFile,
                height: 200,
                width: double.infinity,
                fit: BoxFit.cover,
              ),
              SizedBox(height: 16),
              TextField(
                controller: nameController,
                decoration: InputDecoration(
                  labelText: 'Name *',
                  border: OutlineInputBorder(),
                ),
              ),
              SizedBox(height: 8),
              TextField(
                controller: aadharController,
                decoration: InputDecoration(
                  labelText: 'Aadhar Number *',
                  border: OutlineInputBorder(),
                ),
              ),
              SizedBox(height: 8),
              TextField(
                controller: contactController,
                decoration: InputDecoration(
                  labelText: 'Contact Number',
                  border: OutlineInputBorder(),
                ),
              ),
              SizedBox(height: 8),
              TextField(
                controller: placeController,
                decoration: InputDecoration(
                  labelText: 'Place Lost *',
                  border: OutlineInputBorder(),
                ),
              ),
              SizedBox(height: 8),
              TextField(
                controller: addressController,
                decoration: InputDecoration(
                  labelText: 'Permanent Address *',
                  border: OutlineInputBorder(),
                ),
                maxLines: 3,
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              if (nameController.text.isEmpty ||
                  aadharController.text.isEmpty ||
                  placeController.text.isEmpty ||
                  addressController.text.isEmpty) {
                _showErrorSnackBar('Please fill all required fields');
                return;
              }

              Navigator.pop(context);

              setState(() {
                isLoading = true;
              });

              var result = await LostPersonService.createLostPersonReport(
                name: nameController.text,
                aadharNumber: aadharController.text,
                contactNumber: contactController.text.isNotEmpty
                    ? contactController.text
                    : null,
                placeLost: placeController.text,
                permanentAddress: addressController.text,
                imageFile: imageFile,
              );

              setState(() {
                isLoading = false;
              });

              if (result['success']) {
                _showSuccessSnackBar('Report created successfully');
                _loadLostPersons(); // Reload the list
              } else {
                _showErrorSnackBar(result['error']);
              }
            },
            child: Text('Submit'),
          ),
        ],
      ),
    );
  }

  void _searchLostPersons() {
    final searchController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Search Lost Persons'),
        content: TextField(
          controller: searchController,
          decoration: InputDecoration(
            labelText: 'Search by name or Aadhar number',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              if (searchController.text.isEmpty) {
                _showErrorSnackBar('Please enter a search term');
                return;
              }

              Navigator.pop(context);

              setState(() {
                isLoading = true;
              });

              var result = await LostPersonService.searchLostPersons(
                searchController.text,
              );

              setState(() {
                isLoading = false;
              });

              if (result['success']) {
                var lostPersonList = LostPersonList.fromJson(result['data']);
                setState(() {
                  lostPersons = lostPersonList.lostPersons;
                });
                _showSuccessSnackBar('Search completed');
              } else {
                _showErrorSnackBar(result['error']);
              }
            },
            child: Text('Search'),
          ),
        ],
      ),
    );
  }

  void _showErrorSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
      ),
    );
  }

  void _showSuccessSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.green,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Lost Person Reports'),
        actions: [
          IconButton(
            icon: Icon(Icons.search),
            onPressed: _searchLostPersons,
          ),
        ],
      ),
      body: isLoading
          ? Center(child: CircularProgressIndicator())
          : lostPersons.isEmpty
              ? Center(
                  child: Text(
                    'No lost person reports found',
                    style: TextStyle(fontSize: 18),
                  ),
                )
              : ListView.builder(
                  itemCount: lostPersons.length,
                  itemBuilder: (context, index) {
                    final person = lostPersons[index];
                    return Card(
                      margin: EdgeInsets.all(8),
                      child: ListTile(
                        leading: CircleAvatar(
                          backgroundImage: NetworkImage(
                            LostPersonService.getImageUrl(person.imagePath),
                          ),
                        ),
                        title: Text(person.name),
                        subtitle: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text('Aadhar: ${person.aadharNumber}'),
                            Text('Place Lost: ${person.placeLost}'),
                            if (person.contactNumber != null)
                              Text('Contact: ${person.contactNumber}'),
                          ],
                        ),
                        trailing: PopupMenuButton(
                          itemBuilder: (context) => [
                            PopupMenuItem(
                              value: 'view',
                              child: Text('View Details'),
                            ),
                            PopupMenuItem(
                              value: 'delete',
                              child: Text('Delete'),
                            ),
                          ],
                          onSelected: (value) async {
                            if (value == 'view') {
                              _showPersonDetails(person);
                            } else if (value == 'delete') {
                              _deletePerson(person.id);
                            }
                          },
                        ),
                      ),
                    );
                  },
                ),
      floatingActionButton: FloatingActionButton(
        onPressed: _addLostPerson,
        child: Icon(Icons.add),
        tooltip: 'Add Lost Person Report',
      ),
    );
  }

  void _showPersonDetails(LostPerson person) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Person Details'),
        content: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Image.network(
                LostPersonService.getImageUrl(person.imagePath),
                height: 200,
                width: double.infinity,
                fit: BoxFit.cover,
              ),
              SizedBox(height: 16),
              Text('Name: ${person.name}'),
              Text('Aadhar Number: ${person.aadharNumber}'),
              if (person.contactNumber != null)
                Text('Contact Number: ${person.contactNumber}'),
              Text('Place Lost: ${person.placeLost}'),
              Text('Permanent Address: ${person.permanentAddress}'),
              Text('Reported on: ${person.uploadTimestamp.toString()}'),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Close'),
          ),
        ],
      ),
    );
  }

  Future<void> _deletePerson(String id) async {
    bool confirm = await showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Confirm Delete'),
        content: Text('Are you sure you want to delete this report?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: Text('Delete'),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
          ),
        ],
      ),
    ) ?? false;

    if (confirm) {
      setState(() {
        isLoading = true;
      });

      var result = await LostPersonService.deleteLostPersonReport(id);

      setState(() {
        isLoading = false;
      });

      if (result['success']) {
        _showSuccessSnackBar('Report deleted successfully');
        _loadLostPersons(); // Reload the list
      } else {
        _showErrorSnackBar(result['error']);
      }
    }
  }
}
```

## Key Integration Points

### 1. **Image Upload**
- Use `image_picker` package for camera/gallery access
- Request camera permissions before accessing camera
- Compress images to reduce upload size
- Handle file selection errors gracefully

### 2. **Form Validation**
- Validate required fields before submission
- Show appropriate error messages
- Handle network errors and timeouts

### 3. **State Management**
- Use `setState` for simple state management
- Consider using Provider, Bloc, or Riverpod for complex apps
- Show loading indicators during API calls

### 4. **Error Handling**
- Handle network errors
- Show user-friendly error messages
- Implement retry mechanisms

### 5. **Image Display**
- Use `NetworkImage` to display uploaded images
- Handle image loading errors
- Implement image caching for better performance

## Important Notes for Frontend Developer

### 1. **Base URL Configuration**
- Change the `baseUrl` in `LostPersonService` to match your backend server
- For development: `http://localhost:8080/api/v1`
- For production: `https://your-domain.com/api/v1`

### 2. **File Upload Requirements**
- Maximum file size: 10MB
- Supported formats: JPEG, JPG, PNG, GIF
- Field name must be `image`

### 3. **Required Fields**
- Name (required)
- Aadhar Number (required)
- Place Lost (required)
- Permanent Address (required)
- Contact Number (optional)
- Image (required)

### 4. **Response Format**
All API responses follow this format:
```json
{
  "success": true/false,
  "message": "Response message",
  "data": { ... },
  "error": "Error message (if success is false)"
}
```

### 5. **Error Handling**
- Always check the `success` field in responses
- Display error messages from the `error` field
- Handle network timeouts and connection errors

### 6. **Image URLs**
- Images are served at `/api/v1/images/{filename}`
- Use the `getImageUrl()` method to construct proper URLs
- Handle cases where images might not load

### 7. **Testing**
- Test with different image sizes and formats
- Test network connectivity scenarios
- Test form validation thoroughly
- Test on both Android and iOS devices

## Troubleshooting

### Common Issues:

1. **Network Error**: Check if backend server is running
2. **Image Upload Fails**: Verify file size and format
3. **Permission Denied**: Request camera/storage permissions
4. **CORS Issues**: Backend should allow requests from your Flutter app
5. **Image Not Loading**: Check image URL construction

### Debug Tips:

1. Use `print()` statements to debug API calls
2. Check network tab in browser dev tools
3. Verify request/response format
4. Test API endpoints with Postman first

This integration guide provides everything needed to connect your Flutter app with the TrinetraGuard backend. The service class handles all API communication, and the example app shows how to use it in a real application.
