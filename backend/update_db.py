#!/usr/bin/env python3
import sqlite3
import json
import os

# Connect to database
conn = sqlite3.connect('database.db')
cursor = conn.cursor()

# Update video with total people count
cursor.execute("UPDATE videos SET total_people = ? WHERE id = ?", (8, 1))

# Clear existing faces for this video
cursor.execute("DELETE FROM faces WHERE video_id = ?", (1,))

# Insert the 8 faces
faces_data = [
    ("person_16.jpg", '["0:00:07"]'),
    ("person_12.jpg", '["0:00:07"]'),
    ("person_18.jpg", '["0:00:07", "0:00:08"]'),
    ("person_10.jpg", '["0:00:06"]'),
    ("person_9.jpg", '["0:00:06"]'),
    ("person_8.jpg", '["0:00:06"]'),
    ("person_14.jpg", '["0:00:07"]'),
    ("person_3.jpg", '["0:00:01"]')
]

for image_path, timestamps in faces_data:
    cursor.execute(
        "INSERT INTO faces (video_id, image_path, timestamps, name) VALUES (?, ?, ?, ?)",
        (1, image_path, timestamps, "Unknown")
    )

# Commit changes
conn.commit()
conn.close()

print("Database updated successfully!")
print("Video ID 1 now has 8 unique faces detected") 