# React Integration Guide for TrinetraGuard Backend

## Quick Start

### 1. Create React App
```bash
npm create vite@latest lost-person-app -- --template react
cd lost-person-app
npm install axios
```

### 2. API Service
```javascript
// src/services/api.js
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export const api = {
  // Get all lost persons
  getAllLostPersons: async () => {
    const response = await axios.get(`${API_BASE_URL}/lost-persons/`);
    return response.data.data;
  },

  // Search lost persons
  searchLostPersons: async (query) => {
    const response = await axios.get(`${API_BASE_URL}/lost-persons/search?q=${encodeURIComponent(query)}`);
    return response.data.data;
  },

  // Get image URL
  getImageUrl: (imagePath) => {
    const filename = imagePath.split('/').pop();
    return `${API_BASE_URL}/images/${filename}`;
  }
};
```

### 3. Main App Component
```jsx
// src/App.jsx
import React, { useState, useEffect } from 'react';
import { api } from './services/api';
import './App.css';

function App() {
  const [lostPersons, setLostPersons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    fetchLostPersons();
  }, []);

  const fetchLostPersons = async () => {
    try {
      setLoading(true);
      const data = await api.getAllLostPersons();
      setLostPersons(data.lost_persons);
    } catch (error) {
      console.error('Error fetching lost persons:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (query) => {
    if (!query.trim()) {
      fetchLostPersons();
      return;
    }

    try {
      setLoading(true);
      const data = await api.searchLostPersons(query);
      setLostPersons(data.lost_persons);
    } catch (error) {
      console.error('Error searching:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="app">
      <header>
        <h1>🔍 Lost Person Reports</h1>
        <div className="search-container">
          <input
            type="text"
            placeholder="Search by name or Aadhar number..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSearch(searchQuery)}
          />
          <button onClick={() => handleSearch(searchQuery)}>Search</button>
          <button onClick={fetchLostPersons}>Clear</button>
        </div>
      </header>

      <main>
        <div className="lost-persons-grid">
          {lostPersons.map((person) => (
            <div key={person.id} className="person-card">
              <div className="person-image">
                <img
                  src={api.getImageUrl(person.image_path)}
                  alt={person.name}
                  onError={(e) => {
                    e.target.style.display = 'none';
                    e.target.nextSibling.style.display = 'flex';
                  }}
                />
                <div className="image-placeholder" style={{ display: 'none' }}>
                  <span>📷</span>
                  <p>Image not available</p>
                </div>
              </div>
              
              <div className="person-info">
                <h3>{person.name}</h3>
                <p><strong>Aadhar:</strong> {person.aadhar_number}</p>
                <p><strong>Place Lost:</strong> {person.place_lost}</p>
                {person.contact_number && (
                  <p><strong>Contact:</strong> {person.contact_number}</p>
                )}
                <p><strong>Address:</strong> {person.permanent_address}</p>
                <p><strong>Reported:</strong> {new Date(person.upload_timestamp).toLocaleDateString()}</p>
              </div>
            </div>
          ))}
        </div>

        {lostPersons.length === 0 && (
          <div className="empty-state">
            <p>No lost person reports found</p>
          </div>
        )}
      </main>
    </div>
  );
}

export default App;
```

### 4. CSS Styles
```css
/* src/App.css */
.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

header {
  text-align: center;
  margin-bottom: 30px;
}

.search-container {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 20px;
}

.search-container input {
  padding: 10px;
  border: 2px solid #ddd;
  border-radius: 5px;
  width: 300px;
}

.search-container button {
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.lost-persons-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.person-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.person-image {
  position: relative;
  height: 200px;
}

.person-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #666;
}

.person-info {
  padding: 15px;
}

.person-info h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.person-info p {
  margin: 5px 0;
  color: #666;
}

.loading {
  text-align: center;
  padding: 50px;
  font-size: 18px;
}

.empty-state {
  text-align: center;
  padding: 50px;
  color: #666;
}
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/lost-persons/` | Get all lost person reports |
| `GET` | `/api/v1/lost-persons/search?q=query` | Search reports |
| `GET` | `/api/v1/images/{filename}` | Get image file |

## Response Format

```json
{
  "success": true,
  "message": "Lost person reports retrieved successfully",
  "data": {
    "total": 2,
    "lost_persons": [
      {
        "id": "uuid",
        "name": "John Doe",
        "aadhar_number": "123456789012",
        "contact_number": "9876543210",
        "place_lost": "Mumbai Central",
        "permanent_address": "123 Main St",
        "image_path": "uploads/filename.jpg",
        "upload_timestamp": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

## Key Features

✅ **View All Reports** - Display all lost person reports in a grid
✅ **Search Functionality** - Search by name or Aadhar number
✅ **Image Display** - Show uploaded images with fallback
✅ **Responsive Design** - Works on all screen sizes
✅ **Error Handling** - Graceful error handling for API calls

## Important Notes

1. **Backend URL**: Change `API_BASE_URL` to match your server
2. **CORS**: Ensure backend allows requests from your React app
3. **Image URLs**: Images are served directly from backend
4. **Error Handling**: Add proper error handling for production

## Testing

1. Start backend: `go run main.go`
2. Start React app: `npm run dev`
3. Test search functionality
4. Verify image display
5. Test responsive design

This provides a complete React frontend for viewing and searching lost person reports with images.
