# React Vite Integration Guide for TrinetraGuard Backend

## Overview

This guide provides step-by-step instructions for integrating the TrinetraGuard Lost Person Backend with your React Vite application. The backend provides RESTful API endpoints for managing lost person reports with image display capabilities.

## Backend Information

- **Base URL**: `http://localhost:8080/api/v1` (change to your server URL)
- **Response Format**: JSON
- **Image URLs**: `http://localhost:8080/api/v1/images/{filename}`
- **Authentication**: None required (for now)

## Prerequisites

1. **Node.js**: Ensure you have Node.js installed
2. **React Vite**: Create a new Vite project or use existing one
3. **Axios**: For HTTP requests
4. **React Router**: For navigation (optional)

## Setup Dependencies

Create a new Vite project or add dependencies to existing project:

```bash
# Create new Vite project
npm create vite@latest lost-person-app -- --template react
cd lost-person-app

# Install dependencies
npm install axios react-router-dom
```

## API Service Class

Create a service class to handle all API calls:

```javascript
// src/services/lostPersonService.js
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

class LostPersonService {
  // Get all lost person reports
  static async getAllLostPersons() {
    try {
      const response = await axios.get(`${API_BASE_URL}/lost-persons/`);
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to fetch lost persons'
      };
    }
  }

  // Get a specific lost person report by ID
  static async getLostPersonById(id) {
    try {
      const response = await axios.get(`${API_BASE_URL}/lost-persons/${id}`);
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to fetch lost person'
      };
    }
  }

  // Search lost persons
  static async searchLostPersons(query) {
    try {
      const response = await axios.get(`${API_BASE_URL}/lost-persons/search?q=${encodeURIComponent(query)}`);
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to search lost persons'
      };
    }
  }

  // Create a new lost person report
  static async createLostPersonReport(formData) {
    try {
      const response = await axios.post(`${API_BASE_URL}/lost-persons/`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to create lost person report'
      };
    }
  }

  // Update a lost person report
  static async updateLostPersonReport(id, formData) {
    try {
      const response = await axios.put(`${API_BASE_URL}/lost-persons/${id}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return {
        success: true,
        data: response.data.data,
        message: response.data.message
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to update lost person report'
      };
    }
  }

  // Delete a lost person report
  static async deleteLostPersonReport(id) {
    try {
      const response = await axios.delete(`${API_BASE_URL}/lost-persons/${id}`);
      return {
        success: true,
        message: response.data.message
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to delete lost person report'
      };
    }
  }

  // Get image URL
  static getImageUrl(imagePath) {
    // Extract filename from path (e.g., "uploads/filename.jpg" -> "filename.jpg")
    const filename = imagePath.split('/').pop();
    return `${API_BASE_URL}/images/${filename}`;
  }
}

export default LostPersonService;
```

## Main App Component

Create the main app component to display all lost persons:

```jsx
// src/App.jsx
import React, { useState, useEffect } from 'react';
import LostPersonService from './services/lostPersonService';
import LostPersonCard from './components/LostPersonCard';
import SearchBar from './components/SearchBar';
import AddPersonForm from './components/AddPersonForm';
import './App.css';

function App() {
  const [lostPersons, setLostPersons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);

  // Fetch all lost persons on component mount
  useEffect(() => {
    fetchLostPersons();
  }, []);

  const fetchLostPersons = async () => {
    setLoading(true);
    setError(null);
    
    const result = await LostPersonService.getAllLostPersons();
    
    if (result.success) {
      setLostPersons(result.data.lost_persons);
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  const handleSearch = async (query) => {
    if (!query.trim()) {
      fetchLostPersons();
      return;
    }

    setLoading(true);
    setError(null);
    
    const result = await LostPersonService.searchLostPersons(query);
    
    if (result.success) {
      setLostPersons(result.data.lost_persons);
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  const handleAddPerson = async (formData) => {
    const result = await LostPersonService.createLostPersonReport(formData);
    
    if (result.success) {
      setShowAddForm(false);
      fetchLostPersons(); // Refresh the list
      alert('Lost person report created successfully!');
    } else {
      alert(`Error: ${result.error}`);
    }
  };

  const handleDeletePerson = async (id) => {
    if (window.confirm('Are you sure you want to delete this report?')) {
      const result = await LostPersonService.deleteLostPersonReport(id);
      
      if (result.success) {
        fetchLostPersons(); // Refresh the list
        alert('Report deleted successfully!');
      } else {
        alert(`Error: ${result.error}`);
      }
    }
  };

  if (loading) {
    return (
      <div className="app">
        <div className="loading">
          <div className="spinner"></div>
          <p>Loading lost person reports...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>🔍 Lost Person Reports</h1>
        <p>Find and manage lost person reports</p>
      </header>

      <main className="app-main">
        <div className="controls">
          <SearchBar onSearch={handleSearch} />
          <button 
            className="add-button"
            onClick={() => setShowAddForm(true)}
          >
            ➕ Add New Report
          </button>
        </div>

        {error && (
          <div className="error-message">
            <p>❌ Error: {error}</p>
            <button onClick={fetchLostPersons}>🔄 Retry</button>
          </div>
        )}

        {lostPersons.length === 0 && !error ? (
          <div className="empty-state">
            <p>📭 No lost person reports found</p>
            <p>Click "Add New Report" to create the first report</p>
          </div>
        ) : (
          <div className="lost-persons-grid">
            {lostPersons.map((person) => (
              <LostPersonCard
                key={person.id}
                person={person}
                onDelete={handleDeletePerson}
              />
            ))}
          </div>
        )}
      </main>

      {showAddForm && (
        <AddPersonForm
          onSubmit={handleAddPerson}
          onCancel={() => setShowAddForm(false)}
        />
      )}
    </div>
  );
}

export default App;
```

## Lost Person Card Component

Create a component to display individual lost person cards:

```jsx
// src/components/LostPersonCard.jsx
import React, { useState } from 'react';
import LostPersonService from '../services/lostPersonService';
import './LostPersonCard.css';

function LostPersonCard({ person, onDelete }) {
  const [showDetails, setShowDetails] = useState(false);
  const [imageError, setImageError] = useState(false);

  const imageUrl = LostPersonService.getImageUrl(person.image_path);

  const handleImageError = () => {
    setImageError(true);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="lost-person-card">
      <div className="card-header">
        <div className="person-image">
          {imageError ? (
            <div className="image-placeholder">
              <span>📷</span>
              <p>Image not available</p>
            </div>
          ) : (
            <img
              src={imageUrl}
              alt={`${person.name}`}
              onError={handleImageError}
              loading="lazy"
            />
          )}
        </div>
        
        <div className="card-actions">
          <button
            className="action-btn details-btn"
            onClick={() => setShowDetails(!showDetails)}
          >
            {showDetails ? '👁️ Hide' : '👁️ Details'}
          </button>
          <button
            className="action-btn delete-btn"
            onClick={() => onDelete(person.id)}
          >
            🗑️ Delete
          </button>
        </div>
      </div>

      <div className="card-content">
        <h3 className="person-name">{person.name}</h3>
        <p className="person-aadhar">Aadhar: {person.aadhar_number}</p>
        <p className="person-place">📍 Lost at: {person.place_lost}</p>
        
        {person.contact_number && (
          <p className="person-contact">📞 Contact: {person.contact_number}</p>
        )}
        
        <p className="person-date">
          📅 Reported: {formatDate(person.upload_timestamp)}
        </p>
      </div>

      {showDetails && (
        <div className="card-details">
          <h4>📋 Full Details</h4>
          <div className="details-grid">
            <div className="detail-item">
              <strong>Name:</strong> {person.name}
            </div>
            <div className="detail-item">
              <strong>Aadhar Number:</strong> {person.aadhar_number}
            </div>
            {person.contact_number && (
              <div className="detail-item">
                <strong>Contact Number:</strong> {person.contact_number}
              </div>
            )}
            <div className="detail-item">
              <strong>Place Lost:</strong> {person.place_lost}
            </div>
            <div className="detail-item full-width">
              <strong>Permanent Address:</strong> {person.permanent_address}
            </div>
            <div className="detail-item">
              <strong>Report ID:</strong> {person.id}
            </div>
            <div className="detail-item">
              <strong>Reported On:</strong> {formatDate(person.upload_timestamp)}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default LostPersonCard;
```

## Search Bar Component

```jsx
// src/components/SearchBar.jsx
import React, { useState } from 'react';
import './SearchBar.css';

function SearchBar({ onSearch }) {
  const [query, setQuery] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSearch(query);
  };

  const handleClear = () => {
    setQuery('');
    onSearch('');
  };

  return (
    <form className="search-bar" onSubmit={handleSubmit}>
      <div className="search-input-group">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search by name or Aadhar number..."
          className="search-input"
        />
        <button type="submit" className="search-btn">
          🔍 Search
        </button>
        {query && (
          <button type="button" onClick={handleClear} className="clear-btn">
            ❌ Clear
          </button>
        )}
      </div>
    </form>
  );
}

export default SearchBar;
```

## Add Person Form Component

```jsx
// src/components/AddPersonForm.jsx
import React, { useState } from 'react';
import './AddPersonForm.css';

function AddPersonForm({ onSubmit, onCancel }) {
  const [formData, setFormData] = useState({
    name: '',
    aadharNumber: '',
    contactNumber: '',
    placeLost: '',
    permanentAddress: ''
  });
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setSelectedImage(file);
      
      // Create preview
      const reader = new FileReader();
      reader.onload = (e) => {
        setImagePreview(e.target.result);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validate required fields
    if (!formData.name || !formData.aadharNumber || !formData.placeLost || 
        !formData.permanentAddress || !selectedImage) {
      alert('Please fill all required fields and select an image');
      return;
    }

    setLoading(true);

    try {
      const data = new FormData();
      data.append('name', formData.name);
      data.append('aadhar_number', formData.aadharNumber);
      data.append('contact_number', formData.contactNumber);
      data.append('place_lost', formData.placeLost);
      data.append('permanent_address', formData.permanentAddress);
      data.append('image', selectedImage);

      await onSubmit(data);
    } catch (error) {
      alert('Error creating report: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h2>➕ Add New Lost Person Report</h2>
          <button onClick={onCancel} className="close-btn">✕</button>
        </div>

        <form onSubmit={handleSubmit} className="add-form">
          <div className="form-group">
            <label htmlFor="name">Name *</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="aadharNumber">Aadhar Number *</label>
            <input
              type="text"
              id="aadharNumber"
              name="aadharNumber"
              value={formData.aadharNumber}
              onChange={handleInputChange}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="contactNumber">Contact Number</label>
            <input
              type="tel"
              id="contactNumber"
              name="contactNumber"
              value={formData.contactNumber}
              onChange={handleInputChange}
            />
          </div>

          <div className="form-group">
            <label htmlFor="placeLost">Place Lost *</label>
            <input
              type="text"
              id="placeLost"
              name="placeLost"
              value={formData.placeLost}
              onChange={handleInputChange}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="permanentAddress">Permanent Address *</label>
            <textarea
              id="permanentAddress"
              name="permanentAddress"
              value={formData.permanentAddress}
              onChange={handleInputChange}
              required
              rows="3"
            />
          </div>

          <div className="form-group">
            <label htmlFor="image">Person's Image *</label>
            <input
              type="file"
              id="image"
              accept="image/*"
              onChange={handleImageChange}
              required
            />
            {imagePreview && (
              <div className="image-preview">
                <img src={imagePreview} alt="Preview" />
              </div>
            )}
          </div>

          <div className="form-actions">
            <button type="button" onClick={onCancel} className="cancel-btn">
              Cancel
            </button>
            <button type="submit" className="submit-btn" disabled={loading}>
              {loading ? 'Creating...' : 'Create Report'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default AddPersonForm;
```

## CSS Styles

### App.css
```css
/* src/App.css */
.app {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.app-header {
  text-align: center;
  color: white;
  margin-bottom: 30px;
}

.app-header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.app-header p {
  font-size: 1.1rem;
  opacity: 0.9;
}

.app-main {
  max-width: 1200px;
  margin: 0 auto;
}

.controls {
  display: flex;
  gap: 20px;
  margin-bottom: 30px;
  align-items: center;
  flex-wrap: wrap;
}

.add-button {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.3s;
}

.add-button:hover {
  background: #45a049;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 50vh;
  color: white;
}

.spinner {
  border: 4px solid rgba(255,255,255,0.3);
  border-top: 4px solid white;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  background: #ffebee;
  color: #c62828;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-message button {
  background: #c62828;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.empty-state {
  text-align: center;
  color: white;
  padding: 60px 20px;
}

.empty-state p {
  font-size: 1.2rem;
  margin: 10px 0;
}

.lost-persons-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

@media (max-width: 768px) {
  .controls {
    flex-direction: column;
    align-items: stretch;
  }
  
  .lost-persons-grid {
    grid-template-columns: 1fr;
  }
}
```

### LostPersonCard.css
```css
/* src/components/LostPersonCard.css */
.lost-person-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
}

.lost-person-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.card-header {
  position: relative;
}

.person-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.person-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.person-image:hover img {
  transform: scale(1.05);
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #666;
}

.image-placeholder span {
  font-size: 3rem;
  margin-bottom: 10px;
}

.card-actions {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  gap: 8px;
}

.action-btn {
  background: rgba(0,0,0,0.7);
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background 0.3s;
}

.action-btn:hover {
  background: rgba(0,0,0,0.9);
}

.delete-btn:hover {
  background: #d32f2f;
}

.card-content {
  padding: 20px;
}

.person-name {
  font-size: 1.4rem;
  margin: 0 0 10px 0;
  color: #333;
}

.person-aadhar,
.person-place,
.person-contact,
.person-date {
  margin: 5px 0;
  color: #666;
  font-size: 0.95rem;
}

.card-details {
  padding: 20px;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
}

.card-details h4 {
  margin: 0 0 15px 0;
  color: #333;
}

.details-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.detail-item {
  padding: 8px 0;
}

.detail-item.full-width {
  grid-column: 1 / -1;
}

.detail-item strong {
  color: #333;
  margin-right: 5px;
}
```

### SearchBar.css
```css
/* src/components/SearchBar.css */
.search-bar {
  flex: 1;
  max-width: 500px;
}

.search-input-group {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.search-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.3s;
}

.search-btn:hover {
  background: #5a6fd8;
}

.clear-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.3s;
}

.clear-btn:hover {
  background: #c82333;
}
```

### AddPersonForm.css
```css
/* src/components/AddPersonForm.css */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
}

.modal-header h2 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}

.close-btn:hover {
  color: #333;
}

.add-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input[type="file"] {
  padding: 8px;
}

.image-preview {
  margin-top: 10px;
  text-align: center;
}

.image-preview img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.form-actions {
  display: flex;
  gap: 15px;
  justify-content: flex-end;
  margin-top: 30px;
}

.cancel-btn,
.submit-btn {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.3s;
}

.cancel-btn {
  background: #6c757d;
  color: white;
  border: none;
}

.cancel-btn:hover {
  background: #5a6268;
}

.submit-btn {
  background: #4CAF50;
  color: white;
  border: none;
}

.submit-btn:hover:not(:disabled) {
  background: #45a049;
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
```

## Usage Instructions

1. **Start the Backend Server:**
   ```bash
   cd Trinetr-backend
   go run main.go
   ```

2. **Start the React Vite App:**
   ```bash
   cd lost-person-app
   npm run dev
   ```

3. **Access the Application:**
   - Backend: `http://localhost:8080`
   - Frontend: `http://localhost:5173`

## Key Features

### ✅ **Complete CRUD Operations**
- View all lost person reports
- Search by name or Aadhar number
- Add new reports with image upload
- Delete existing reports
- View detailed information

### ✅ **Image Display**
- Automatic image loading from backend
- Fallback for missing images
- Image preview in add form
- Responsive image grid

### ✅ **Responsive Design**
- Mobile-friendly layout
- Grid adapts to screen size
- Touch-friendly buttons

### ✅ **Error Handling**
- Network error handling
- Form validation
- User-friendly error messages
- Loading states

### ✅ **Search Functionality**
- Real-time search
- Clear search option
- Search by name or Aadhar

## API Endpoints Used

1. **GET** `/api/v1/lost-persons/` - Fetch all reports
2. **GET** `/api/v1/lost-persons/search?q=query` - Search reports
3. **POST** `/api/v1/lost-persons/` - Create new report
4. **DELETE** `/api/v1/lost-persons/:id` - Delete report
5. **GET** `/api/v1/images/:filename` - Get image file

## Important Notes

1. **CORS Configuration**: Ensure your backend allows requests from `http://localhost:5173`
2. **Image URLs**: Images are served directly from the backend
3. **File Upload**: Supports multipart form data for image uploads
4. **Error Handling**: Comprehensive error handling for all API calls
5. **Responsive**: Works on all device sizes

This React Vite integration provides a complete frontend solution for managing lost person reports with full CRUD functionality and image display capabilities.
