package main

import (
	"backend/cmd/api/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// data type of the post stored in the database
type Post struct {
    ID        int    `json:"id"`
    UserID    int    `json:"user_id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
}

// data type of users stored in the database
type User struct {
    ID           int       `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

// data type received from the front end
type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// data type of JWT token
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from %s", app.Domain)
}

// API endpoint for getting posts
func (app *application) GetPosts(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GetPosts endpoint hit")
    
    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Println("Database connection error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // updated query to match the schema
    rows, err := db.Query("SELECT id, user_id, title, content, created_at FROM posts")
    if err != nil {
        fmt.Println("Query error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
        if err != nil {
            fmt.Println("Row scan error:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

	// print if there is any error
    if err := json.NewEncoder(w).Encode(posts); err != nil {
        fmt.Println("JSON encoding error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Printf("Returning %d posts\n", len(posts))
}

// API endpoint for logging in
func (app *application) Login(w http.ResponseWriter, r *http.Request) {

// parsing the incoming JSON request and assigning it to the creds variable
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
				fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}
		fmt.Printf("Received credentials - Email: %s, Password: %s\n", creds.Email, creds.Password)

	// connect to the database
	connStr := config.GetDBConfig()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
				fmt.Printf("Database connection error: %v\n", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
		fmt.Println("Database connection successful")

	// fetch user from the database and put then into our user variable
	var user User
	err = db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", creds.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
				fmt.Printf("User lookup error: %v\n", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
		fmt.Printf("Found user with ID: %d, Email: %s\n", user.ID, user.Email)

	// Compare passwords by converting both to a byte slice
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
				fmt.Printf("Password comparison failed: %v\n", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return  // Added missing return statement
	}
		fmt.Println("Password verification successful")

	/* generate JWT token, claims are customised information while
	Registered claims are standardized information */
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims {
		UserID: user.ID, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		}, 
	}

	// generation of token with secret key to prevent tampering of token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.JwtSecret))
	if err != nil {
				fmt.Printf("Token generation error: %v\n", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
		fmt.Println("JWT token generated successfully")

	// send response 
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString, 
	})
}