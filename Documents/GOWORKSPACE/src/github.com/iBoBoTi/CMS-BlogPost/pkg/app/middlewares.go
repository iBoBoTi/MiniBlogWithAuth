package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckLoginMiddleware(c *gin.Context){
	session, err := c.Cookie("session")
	if err != nil{
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.Set("sessionToken", session)
	c.Next()
}

func CheckNotLoginMiddleware(c *gin.Context) {
	_, err := c.Cookie("session")
	if err != nil{
		c.Next()
		return
	}
	c.Redirect(http.StatusFound,"/blogar")
	return
}

