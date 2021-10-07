package app

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/repository"
	"log"
	"net/http"
	"strings"
	"github.com/google/uuid"
)

var Dbase repository.Storage

func init(){
	db,_ := repository.DataBaseConnection()
	Dbase.DB = db
}

func index(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

func home(c *gin.Context){
	c.HTML(http.StatusOK, "home.html", nil)
}

func handleSignUp(c *gin.Context){
	c.HTML(http.StatusOK, "SignUp.html", nil)
}

func handleSignUpAuth(c *gin.Context){
	// function will handle the form submitted at signup and add the details to the database
	// validate form takes in a form type return bool true results in processing the details while false results
	// in redirecting to the error string
	checkString,check := ValidateSignUpForm(c)

	if check == true{
		// process form and add to database

		firstname := c.PostForm("FirstName")
		lastname := c.PostForm("LastName")
		username := strings.Trim(c.PostForm("UserName")," ")
		email := c.PostForm("Email")
		password,_ := HashPassword(c.PostForm("password"))

		stmt:= "INSERT INTO `blog-cms`.`users` (`first_name`, `last_name`, `user_name`,`email`,`password`) VALUES (?,?,?,?,?);"
		prepare, err := Dbase.DB.Prepare(stmt)
		errCheck(err)
		defer prepare.Close()
		_,err = prepare.Exec(firstname,lastname,username,email,password)
		if err != nil {
			log.Print(err.Error())
		}
		// redirects to login
		c.Redirect(http.StatusFound,"/login")


	}else{
		c.String(http.StatusBadRequest,checkString)
	}

}

func handleLogin(c *gin.Context){
	// function will present the login form to the user
	c.HTML(http.StatusOK, "login.html", nil)
}

func handleLoginAuth(c *gin.Context){
	// function will collect details from the login form
	// function will verify details as well as create a session if verification passes
	loginCheck := ValidateLoginForm(c)
	sessionToken := uuid.NewString()
	if loginCheck==true{
		c.SetCookie("session", sessionToken, 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound,"/blogar")
	}else{
		c.String(http.StatusUnauthorized,"Details doesn't exist")
	}
}

func handleLogOut(c *gin.Context){
	// function will delete our session and redirect to the home page
	c.SetCookie("session","",-1,"/","localhost",true,true)
	c.Redirect(http.StatusFound,"/")
}

func handleUsersList(){}

func handleUserProfile(){}

func handlePostCreate(){}

func handlePostUpdate(){}

func handlePostRetrieve(){}

func handlePostDelete(){}

func handlePostList(){}






