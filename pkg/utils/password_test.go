package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	// Test hashing
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if hashedPassword == "" {
		t.Fatal("Hashed password should not be empty")
	}
	
	if hashedPassword == password {
		t.Fatal("Hashed password should not be the same as original password")
	}
	
	// Test verification with correct password
	err = VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Failed to verify correct password: %v", err)
	}
	
	// Test verification with incorrect password
	err = VerifyPassword(hashedPassword, "wrongpassword")
	if err == nil {
		t.Fatal("Should fail to verify incorrect password")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	_, err := HashPassword("")
	if err == nil {
		t.Fatal("Should fail to hash empty password")
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testpassword123"
	
	// Hash the same password twice
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password first time: %v", err)
	}
	
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password second time: %v", err)
	}
	
	// Hashes should be different (due to salt)
	if hash1 == hash2 {
		t.Fatal("Two hashes of the same password should be different due to salt")
	}
	
	// But both should verify correctly
	if err := VerifyPassword(hash1, password); err != nil {
		t.Fatalf("First hash should verify correctly: %v", err)
	}
	
	if err := VerifyPassword(hash2, password); err != nil {
		t.Fatalf("Second hash should verify correctly: %v", err)
	}
}
