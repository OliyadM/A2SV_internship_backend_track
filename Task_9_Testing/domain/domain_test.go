package domain

import (
	"encoding/json"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTaskSerialization(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		task := Task{
			ID:          "1",
			Title:       "Test Task",
			Description: "A test task",
			DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
			Status:      "pending",
		}

		
		jsonData, err := json.Marshal(task)
		assert.NoError(t, err, "Task should marshal to JSON without error")
		assert.Contains(t, string(jsonData), `"id":"1"`, "JSON should include id")
		assert.Contains(t, string(jsonData), `"title":"Test Task"`, "JSON should include title")
		assert.Contains(t, string(jsonData), `"description":"A test task"`, "JSON should include description")
		assert.Contains(t, string(jsonData), `"due_date":"2025-04-03T00:00:00Z"`, "JSON should include due_date")
		assert.Contains(t, string(jsonData), `"status":"pending"`, "JSON should include status")

		
		var unmarshaled Task
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err, "Task should unmarshal from JSON without error")
		assert.Equal(t, task, unmarshaled, "Unmarshaled task should match original")
	})

	t.Run("BSON", func(t *testing.T) {
		task := Task{
			ID:          "1",
			Title:       "Test Task",
			Description: "A test task",
			DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
			Status:      "pending",
		}

		
		bsonData, err := bson.Marshal(task)
		assert.NoError(t, err, "Task should marshal to BSON without error")

		
		var unmarshaled Task
		err = bson.Unmarshal(bsonData, &unmarshaled)
		assert.NoError(t, err, "Task should unmarshal from BSON without error")
		assert.Equal(t, task, unmarshaled, "Unmarshaled task should match original")
	})
}

func TestUserSerialization(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		user := User{
			ID:       "1",
			Username: "testuser",
			Password: "hashedpass",
			Role:     "user",
		}

		
		jsonData, err := json.Marshal(user)
		assert.NoError(t, err, "User should marshal to JSON without error")
		assert.Contains(t, string(jsonData), `"id":"1"`, "JSON should include id")
		assert.Contains(t, string(jsonData), `"username":"testuser"`, "JSON should include username")
		assert.Contains(t, string(jsonData), `"password":"hashedpass"`, "JSON should include password")
		assert.Contains(t, string(jsonData), `"role":"user"`, "JSON should include role")

		
		var unmarshaled User
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err, "User should unmarshal from JSON without error")
		assert.Equal(t, user, unmarshaled, "Unmarshaled user should match original")
	})

	t.Run("BSON", func(t *testing.T) {
		user := User{
			ID:       "1",
			Username: "testuser",
			Password: "hashedpass",
			Role:     "user",
		}

		
		bsonData, err := bson.Marshal(user)
		assert.NoError(t, err, "User should marshal to BSON without error")

	
		var unmarshaled User
		err = bson.Unmarshal(bsonData, &unmarshaled)
		assert.NoError(t, err, "User should unmarshal from BSON without error")
		assert.Equal(t, user, unmarshaled, "Unmarshaled user should match original")
	})
}