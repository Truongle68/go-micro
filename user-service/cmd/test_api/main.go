package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

// Response envelopes
type ResponseEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type VerifyOTPResponse struct {
	VerifyOTPToken string `json:"verify_otp_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

type ProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FullName  string    `json:"full_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ForgotPasswordResponse struct {
	ResetToken string `json:"reset_token"`
}

const (
	baseURL   = "http://localhost:4000/api/v1"
	dbURL     = "postgresql://postgres:postgres@localhost:5433/identity-db?sslmode=disable"
	redisAddr = "localhost:6379"
	redisPass = "secretredispass"
	redisDB   = 1

	testPhone    = "+84999999999"
	testUsername = "testuser_1"
	testPassword = "password123"
	testFullName = "Test User"
	testEmail    = "testuser_1@example.com"
)

func main() {
	fmt.Println("=== Starting API Integration Tests ===")

	// 1. Cleanup database
	fmt.Println("[Step 1] Connecting to Database and cleaning up existing test users...")
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		fmt.Printf("FAIL: sql.Open error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("FAIL: DB Ping error: %v\n", err)
		os.Exit(1)
	}

	// Delete from tables in correct dependency order by querying IDs
	var userID string
	err = db.QueryRowContext(ctx, "SELECT user_id FROM user_credentials WHERE identifier = $1", testPhone).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("FAIL: Querying existing user: %v\n", err)
		os.Exit(1)
	}

	if userID != "" {
		_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
		if err != nil {
			fmt.Printf("FAIL: Cleaning up user by ID %s: %v\n", userID, err)
			os.Exit(1)
		}
		fmt.Printf("Cleaned up user ID: %s\n", userID)
	}

	// Delete user by username if it exists
	var userIDByUsername string
	err = db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", testUsername).Scan(&userIDByUsername)
	if err == nil && userIDByUsername != "" {
		_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userIDByUsername)
		if err != nil {
			fmt.Printf("FAIL: Cleaning up user by username %s: %v\n", testUsername, err)
			os.Exit(1)
		}
		fmt.Printf("Cleaned up user by username ID: %s\n", userIDByUsername)
	}

	// 2. Connect to Redis to fetch OTPs later
	fmt.Println("[Step 2] Connecting to Redis...")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDB,
	})
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("FAIL: Redis Ping error: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// 3. Request OTP (Register)
	fmt.Println("[Step 3] POST /auth/register...")
	reqBody, _ := json.Marshal(map[string]string{
		"phone": testPhone,
	})
	resp, err := client.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: register request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	// 4. Retrieve OTP from Redis
	fmt.Println("[Step 4] Retrieving OTP from Redis...")
	redisKey := fmt.Sprintf("otp:register:%s", testPhone)
	otpCode, err := rdb.Get(ctx, redisKey).Result()
	if err != nil {
		fmt.Printf("FAIL: Could not retrieve OTP code from Redis: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved OTP Code: %s\n", otpCode)

	// 5. Verify OTP
	fmt.Println("[Step 5] POST /auth/register/verify...")
	reqBody, _ = json.Marshal(map[string]string{
		"phone":    testPhone,
		"otp_code": otpCode,
	})
	resp, err = client.Post(baseURL+"/auth/register/verify", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: register/verify request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	var envelope ResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Printf("FAIL: Decode verify-otp envelope failed: %v\n", err)
		os.Exit(1)
	}

	var verifyRes VerifyOTPResponse
	if err := json.Unmarshal(envelope.Data, &verifyRes); err != nil {
		fmt.Printf("FAIL: Unmarshal verify-otp data failed: %v. Raw data: %s\n", err, string(envelope.Data))
		os.Exit(1)
	}
	fmt.Printf("Retrieved Verification Token: %s\n", verifyRes.VerifyOTPToken)

	// 6. Complete Registration
	fmt.Println("[Step 6] POST /auth/register/complete...")
	reqBody, _ = json.Marshal(map[string]string{
		"token":              verifyRes.VerifyOTPToken,
		"full_name":          testFullName,
		"username":           testUsername,
		"password":           testPassword,
		"confirmed_password": testPassword,
	})
	resp, err = client.Post(baseURL+"/auth/register/complete", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: complete-register request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Printf("FAIL: Decode complete-register envelope: %v\n", err)
		os.Exit(1)
	}
	var authRes AuthResponse
	if err := json.Unmarshal(envelope.Data, &authRes); err != nil {
		fmt.Printf("FAIL: Unmarshal complete-register data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Registered successfully. UserID: %s, AccessToken: (present), RefreshToken: (present)\n", authRes.UserID)

	// 7. Get Profile (Success case)
	fmt.Println("[Step 7] GET /users/get-profile...")
	getReq, _ := http.NewRequest("GET", baseURL+"/users/get-profile", nil)
	getReq.Header.Set("Authorization", "Bearer "+authRes.AccessToken)
	resp, err = client.Do(getReq)
	if err != nil {
		fmt.Printf("FAIL: get-profile request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Printf("FAIL: Decode get-profile envelope: %v\n", err)
		os.Exit(1)
	}
	var profileRes ProfileResponse
	if err := json.Unmarshal(envelope.Data, &profileRes); err != nil {
		fmt.Printf("FAIL: Unmarshal get-profile data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Profile details: Username=%s, Phone=%s, FullName=%s\n", profileRes.Username, profileRes.Phone, profileRes.FullName)

	// 8. Update Profile
	fmt.Println("[Step 8] POST /users/update-profile...")
	reqBody, _ = json.Marshal(map[string]string{
		"full_name": "Updated Test User",
	})
	postReq, _ := http.NewRequest("POST", baseURL+"/users/update-profile", bytes.NewBuffer(reqBody))
	postReq.Header.Set("Authorization", "Bearer "+authRes.AccessToken)
	postReq.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(postReq)
	if err != nil {
		fmt.Printf("FAIL: update-profile request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Printf("FAIL: Decode update-profile envelope: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(envelope.Data, &profileRes); err != nil {
		fmt.Printf("FAIL: Unmarshal update-profile data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated Profile details: Phone=%s, FullName=%s\n", profileRes.Phone, profileRes.FullName)

	// 9. Link email to DB for Forgot Password testing
	fmt.Println("[Step 9] Directly injecting email credential in PostgreSQL to test Forgot/Reset Password flow...")
	_, err = db.ExecContext(ctx, `
		INSERT INTO user_credentials (user_id, type, identifier, is_verified, is_primary, secret_hash)
		VALUES ($1, 'email', $2, true, false, '')
	`, authRes.UserID, testEmail)
	if err != nil {
		fmt.Printf("FAIL: Could not inject email credential: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Injected email credential successfully.")

	// Query to list all credentials for debugging
	rows, qErr := db.QueryContext(ctx, "SELECT id, user_id, type, identifier, is_verified, is_primary FROM user_credentials")
	if qErr != nil {
		fmt.Printf("FAIL to query credentials: %v\n", qErr)
	} else {
		defer rows.Close()
		fmt.Println("Existing User Credentials in DB:")
		for rows.Next() {
			var id, uid, t, ident string
			var ver, prim bool
			if err := rows.Scan(&id, &uid, &t, &ident, &ver, &prim); err != nil {
				fmt.Printf("Scan error: %v\n", err)
			} else {
				fmt.Printf(" - ID=%s, UserID=%s, Type=%s, Identifier=%s, Verified=%t, Primary=%t\n", id, uid, t, ident, ver, prim)
			}
		}
	}

	// 10. Forgot Password
	fmt.Println("[Step 10] POST /auth/forgot-password...")
	reqBody, _ = json.Marshal(map[string]string{
		"email": testEmail,
	})
	resp, err = client.Post(baseURL+"/auth/forgot-password", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: forgot-password request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Printf("FAIL: Decode forgot-password envelope: %v\n", err)
		os.Exit(1)
	}
	var forgotRes ForgotPasswordResponse
	if err := json.Unmarshal(envelope.Data, &forgotRes); err != nil {
		fmt.Printf("FAIL: Unmarshal forgot-password data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved Reset Token: %s\n", forgotRes.ResetToken)

	// 11. Reset Password
	fmt.Println("[Step 11] POST /auth/reset-password...")
	newPassword := "newpassword123"
	reqBody, _ = json.Marshal(map[string]string{
		"token":              forgotRes.ResetToken,
		"new_password":       newPassword,
		"confirmed_password": newPassword,
	})
	resp, err = client.Post(baseURL+"/auth/reset-password", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: reset-password request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)
	fmt.Println("Password reset successfully.")

	// 12. Login with old password (should fail)
	fmt.Println("[Step 12] POST /auth/login with OLD password (should fail)...")
	reqBody, _ = json.Marshal(map[string]string{
		"identifier": testUsername,
		"password":   testPassword,
	})
	resp, err = client.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: login request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		fmt.Printf("FAIL: Expected 401 Unauthorized, got status: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	fmt.Println("Correctly failed to log in with old password.")

	// 13. Login with new password (should succeed)
	fmt.Println("[Step 13] POST /auth/login with NEW password...")
	reqBody, _ = json.Marshal(map[string]string{
		"identifier": testUsername,
		"password":   newPassword,
	})
	resp, err = client.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("FAIL: login request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Printf("FAIL: Decode login envelope: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(envelope.Data, &authRes); err != nil {
		fmt.Printf("FAIL: Unmarshal login data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Logged in successfully with new password. UserID: %s\n", authRes.UserID)

	// 14. Logout
	fmt.Println("[Step 14] POST /auth/logout...")
	reqBody, _ = json.Marshal(map[string]string{
		"refresh_token": authRes.RefreshToken,
	})
	logoutReq, _ := http.NewRequest("POST", baseURL+"/auth/logout", bytes.NewBuffer(reqBody))
	logoutReq.Header.Set("Authorization", "Bearer "+authRes.AccessToken)
	logoutReq.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(logoutReq)
	if err != nil {
		fmt.Printf("FAIL: logout request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	verifyStatus(resp, http.StatusOK)
	fmt.Println("Logged out successfully.")

	// 15. Verify token is blacklisted
	fmt.Println("[Step 15] Verifying access token is blacklisted (GET /users/get-profile should fail)...")
	getReq, _ = http.NewRequest("GET", baseURL+"/users/get-profile", nil)
	getReq.Header.Set("Authorization", "Bearer "+authRes.AccessToken)
	resp, err = client.Do(getReq)
	if err != nil {
		fmt.Printf("FAIL: get-profile request after logout failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		fmt.Printf("FAIL: Expected 401 Unauthorized for blacklisted token, got status: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	fmt.Println("Access token successfully verified as blacklisted!")

	fmt.Println("\n=== ALL TESTS COMPLETED SUCCESSFULLY ===")
}

func verifyStatus(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("FAIL: Expected status %d, got %d. Body: %s\n", expected, resp.StatusCode, string(body))
		os.Exit(1)
	}
}
