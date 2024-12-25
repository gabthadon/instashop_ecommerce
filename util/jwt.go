package util

import (
    "github.com/joho/godotenv"
    "github.com/golang-jwt/jwt/v5"
    "log"
    "os"
    "time"
    "errors"
)

// Load environment variables from the .env file
func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }
}

// Use the environment variable for the JWT secret key
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// JWTClaims defines the structure of the JWT claims
type JWTClaims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token
func GenerateJWT(username, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := JWTClaims{
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }
    return signedToken, nil
}

// ValidateJWT parses and validates a JWT token
func ValidateJWT(tokenStr string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, errors.New("token has expired")
        }
        if errors.Is(err, jwt.ErrTokenMalformed) {
            return nil, errors.New("malformed token")
        }
        return nil, errors.New("invalid token")
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
