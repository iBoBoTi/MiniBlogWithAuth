package app

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/api"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Server) ValidateSignUpForm(c *gin.Context) (string, bool) {
	email := c.PostForm("Email")
	username := strings.Trim(c.PostForm("UserName"), " ")
	password := c.PostForm("password")
	if !strings.Contains(email, "@") && !strings.Contains(email, ".") {
		return "check your email: as email should contain \"@\" and \".\"", false
	}
	if username == "" {
		return "please set a username", false
	}

	if len(password) < 4 {
		return "please the length of password should be more than 4", false
	}

	user := api.User{}

	// check the database for already existing username or email
	err := s.DB.DB.QueryRow("SELECT email, user_name FROM users WHERE email=$1 OR user_name=$2 ;", email, username).Scan(&user.Email, &user.UserName)
	errCheck(err)
	if username == user.UserName {
		return "UserName already exist please pick another username", false
	} else if email == user.Email {
		return "Email already exist please use another email to signup", false
	}
	// access database to check for email and username
	return "", true
}

func (s *Server) ValidateLoginForm(c *gin.Context) (bool, string) {
	email := c.PostForm("Email")
	password := c.PostForm("password")
	var user api.User

	err := s.DB.DB.QueryRow(`SELECT "password","id" FROM "users" WHERE email=$1`, email).Scan(&user.Password, &user.ID)
	errCheck(err)
	passCheck := CheckPasswordHash(password, user.Password)
	if passCheck == true {
		return true, user.ID
	}

	return false, ""
}

func errCheck(err error) {
	log.Println(err)
}
