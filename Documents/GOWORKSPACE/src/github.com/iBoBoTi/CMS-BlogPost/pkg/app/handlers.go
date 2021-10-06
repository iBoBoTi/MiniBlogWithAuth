package app

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/repository"
	"net/http"
)



func init(){
	var dbase repository.Storage
	db,_ := repository.DataBaseConnection()
	defer db.Close()

	dbase.DB = db
}

func handleSignUp(c *gin.Context){
	c.HTML(http.StatusOK, "SignUp.html", nil)
}

func handleSignUpAuth(c *gin.Context){

}

func handleLogin(c *gin.Context){
	c.HTML(http.StatusOK, "login.html", nil)
}

func handleLoginAuth(c *gin.Context){}

func handleLogOut(){}

func handleUsersList(){}

func handleUserProfile(){}

func handlePostCreate(){}

func handlePostUpdate(){}

func handlePostRetrieve(){}

func handlePostDelete(){}

func handlePostList(){}

