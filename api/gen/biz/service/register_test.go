package service

import (
	"DouyinMerchant/api/gen/biz/model"
	user "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	return db
}

func TestRegisterService_Run(t *testing.T) {
	//testDB := setupTestDB(t)
	//mysql.DB = testDB // Replace the default DB with the test DB

	ctx := context.Background()

	t.Run("Validation Error - Empty Email", func(t *testing.T) {
		req := &user.RegisterReq{
			Email:           "",
			Password:        "password",
			ConfirmPassword: "password",
		}

		s := NewRegisterService(ctx)
		_, err := s.Run(req)

		if err == nil {
			t.Error("Expected error for empty email, got nil")
		} else if err.Error() != "email is required" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("Validation Error - Password Mismatch", func(t *testing.T) {
		req := &user.RegisterReq{
			Email:           "test@example.com",
			Password:        "password",
			ConfirmPassword: "mismatch",
		}

		s := NewRegisterService(ctx)
		_, err := s.Run(req)

		if err == nil {
			t.Error("Expected error for password mismatch, got nil")
		} else if err.Error() != "passwords do not match" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	//t.Run("Existing Email", func(t *testing.T) {
	//	// Create existing user
	//	existingUser := &mysql.User{
	//		User: &model.User{
	//			Email:    "33@bilibili.com",
	//			Passowrd: "33good",
	//		},
	//	}
	//	if err := mysql.CreateUser(ctx, mysql.DB, existingUser); err != nil {
	//		t.Fatalf("Failed to create existing user: %v", err)
	//	}
	//
	//	req := &user.RegisterReq{
	//		Email:           "existing@example.com",
	//		Password:        "newpassword",
	//		ConfirmPassword: "newpassword",
	//	}
	//
	//	s := NewRegisterService(ctx)
	//	_, err := s.Run(req)
	//
	//	if err == nil {
	//		t.Error("Expected error for existing email, got nil")
	//	} else if err.Error() != "email already exists" {
	//		t.Errorf("Unexpected error message: %v", err)
	//	}
	//})
	//
	//t.Run("Successful Registration", func(t *testing.T) {
	//	req := &user.RegisterReq{
	//		Email:           "newuser@example.com",
	//		Password:        "securepassword",
	//		ConfirmPassword: "securepassword",
	//	}
	//
	//	s := NewRegisterService(ctx)
	//	resp, err := s.Run(req)
	//
	//	if err != nil {
	//		t.Fatalf("Unexpected error: %v", err)
	//	}
	//
	//	if resp.UserId <= 0 {
	//		t.Errorf("Invalid user ID: %v", resp.UserId)
	//	}
	//
	//	// Verify user in the database
	//	var user model.User
	//	if err := mysql.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
	//		t.Fatalf("Failed to find user: %v", err)
	//	}
	//
	//	// Check password hashing
	//	expectedHash := hashPassword(req.Password)
	//	if user.Passowrd != expectedHash {
	//		t.Errorf("Password hash mismatch: got %v, want %v", user.Passowrd, expectedHash)
	//	}
	//})
}
