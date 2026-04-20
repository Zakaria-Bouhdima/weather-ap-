package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abohmeed/auth/authdb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/go-sql-driver/mysql"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/crypto/bcrypt"
)

var secretkey = os.Getenv("JWT_SECRET")

var dbHost string
var dbRoot = "root"
var dbPassword string

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

func main() {
	if os.Getenv("DB_HOST") != "" {
		dbHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		dbPassword = os.Getenv("DB_PASSWORD")
	}
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	authdb.CreateDB(db)
	authdb.CreateTables(db)
	router := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	router.Use(cors.New(corsConfig))
	router.GET("/health", health)
	router.GET("/", health)
	router.POST("/users/:id", loginUser)
	router.POST("/users", createUser)
	router.Run(":8080")
}

type UserCreds struct {
	Username string `json:"user_name"`
	Password string `json:"user_password"`
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": "The auth is running"})
}

func loginUser(c *gin.Context) {
	var uc UserCreds
	err := c.BindJSON(&uc)
	if err != nil {
		fmt.Println("Received invalid JSON for user login")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect or invalid JSON"})
		return
	}
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	u, err := authdb.GetUserByName(uc.Username, db)
	if err != nil {
		fmt.Println(err)
	}
	if u == (authdb.User{}) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bad credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(uc.Password)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bad credentials"})
		return
	}
	token, err := GenerateJWT(u.Name)
	if err != nil {
		fmt.Printf("Error while generating the token: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"JWT": token})
}

func createUser(c *gin.Context) {
	var u authdb.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect or invalid JSON"})
		return
	}
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	result, err := authdb.CreateUser(db, u)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding the user. Please check the logs"})
		return
	} else if !result {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User already exists"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "User added successfully"})
	}
}

func GenerateJWT(userName string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = userName
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
