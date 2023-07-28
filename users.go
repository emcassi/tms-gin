package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type Claim struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GetAllUsers(c *gin.Context) {
	var users []User
	DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user User
	DB.First(&user, c.Param("id"))
	c.JSON(http.StatusOK, user)
}

func GetUserByEmail(email string) (User, error) {
	var user User
	if err := DB.First(&user, "email = ?", email).Error; err != nil {
		fmt.Println(err)
		if err == gorm.ErrRecordNotFound {
       		return User{}, errors.New("User not found")
    	} else {
    		return User{}, errors.New("Failed to fetch user")
    	}
	}
	return user, nil
}

func CreateUser(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for the presence of 'code' and 'price' fields
	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Email, and Password fields are required"})
		return
	}

	if !IsEmailValid(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	if !IsEmailUnique(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
		return
	}

	if !IsValidPassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number"})
		return
	}

	// Hash the password before saving
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Store the hashed password in the user struct
	user.Password = hashedPassword

	DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user User
	result := DB.Delete(&user, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "User deleted"})
}

// JWT Auth
var secretKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateToken(userID uint, email string) (string, error) {
	claims := Claim{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365).Unix(), // Set token to expire in 1 year
			Issuer:    "TMS",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// AuthMiddleware is a middleware function to handle JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.ParseWithClaims(authHeader, &Claim{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(*Claim); ok && token.Valid {
			// Store the user data in the context for future use
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}

func Login(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email (assuming you have a function to retrieve the user from the DB by email)
	// Replace this with your DB query to fetch the user based on the email
	foundUser, err := GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	// Verify the password using bcrypt's CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Password is correct, user is authenticated
	token, err := GenerateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error generating token"})		
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// Helpers

func IsEmailUnique(email string) bool {
	var count int64
	DB.Model(&User{}).Where("email = ?", email).Count(&count)
	return count == 0
}

func IsEmailValid(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}

func IsValidPassword(password string) bool {
	// Check for at least one lowercase letter
	lowercaseRegex := regexp.MustCompile("[a-z]")
	if !lowercaseRegex.MatchString(password) {
		return false
	}

	// Check for at least one uppercase letter
	uppercaseRegex := regexp.MustCompile("[A-Z]")
	if !uppercaseRegex.MatchString(password) {
		return false
	}

	// Check for at least one digit
	digitRegex := regexp.MustCompile("\\d")
	if !digitRegex.MatchString(password) {
		return false
	}

	// Check for the minimum length of 8 characters
	if len(password) < 8 {
		return false
	}

	return true
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
