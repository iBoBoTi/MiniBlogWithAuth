package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/api"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/repository"
	"log"
	"net/http"
	"strings"
	//"github.com/iBoBoTi/CMS-BlogPost/pkg/api"
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
	stmt := "SELECT posts.id, posts.title, posts.content,posts.post_type, users.user_name FROM posts INNER JOIN users ON posts.author = users.id WHERE posts.post_type = 'Public'"

	rows, err := Dbase.DB.Query(stmt)
	if err != nil{
		return
	}
	defer rows.Close()

	var posts []api.Post

	for rows.Next(){
		var p api.Post
		err:= rows.Scan(&p.ID, &p.Title,&p.Content, &p.PostType,&p.UserName)

		if err != nil{
			return
		}
		posts = append(posts, p)

	}
	c.HTML(http.StatusOK, "home.html",posts)
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
		id := uuid.NewString()
		firstname := c.PostForm("FirstName")
		lastname := c.PostForm("LastName")
		username := strings.Trim(c.PostForm("UserName")," ")
		email := c.PostForm("Email")
		password,_ := HashPassword(c.PostForm("password"))

		stmt:= "INSERT INTO `blog-cms`.`users` (`id`,`first_name`, `last_name`, `user_name`,`email`,`password`) VALUES (?,?,?,?,?,?);"
		prepare, err := Dbase.DB.Prepare(stmt)
		errCheck(err)
		defer prepare.Close()
		_,err = prepare.Exec(id,firstname,lastname,username,email,password)
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
	loginCheck, sessionToken := ValidateLoginForm(c)

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

func handlePostCreate(c *gin.Context){
	// presents the form for creating a post
	c.HTML(http.StatusOK, "CreateBlogPostForm.html", nil)
}

func handlePostCreateForm(c *gin.Context){
	// processes the creat blog post form and adds it to the data
	id := uuid.NewString()
	title := c.PostForm("Title")
	content := c.PostForm("Content")
	posttype := c.PostForm("PostType")
	username,_ := c.Cookie("session")


	stmt:= "INSERT INTO `blog-cms`.`posts` (`id`,`author`,`title`, `content`, `post_type`) VALUES (?,?,?,?,?);"
	if title == "" || content == "" {
		c.String(http.StatusUnauthorized,"Please check your entry dear Blogar")
	} else {
		prepare, err := Dbase.DB.Prepare(stmt)
		errCheck(err)
		defer prepare.Close()
		_,err = prepare.Exec(id,username,title,content,posttype)
		if err != nil {
			log.Print(err.Error())
		}
		c.Redirect(http.StatusFound,"/blogar/")
	}

}

func handlePostEdit(c *gin.Context){
	id := c.Param("id")
	row := Dbase.DB.QueryRow("SELECT `title`,`content`,`post_type`,`id` FROM `blog-cms`.`posts` WHERE id=?;",id)
	var post api.Post
	err := row.Scan(&post.Title, &post.Content, &post.PostType,&post.ID)
	errCheck(err)

	c.HTML(http.StatusOK, "EditBlogPostForm.html", post)
}

func handlePostEditForm(c *gin.Context){
	id := c.Param("id")

	title := c.PostForm("Title")
	content := c.PostForm("Content")
	posttype := c.PostForm("PostType")

	stmt, err:= Dbase.DB.Prepare("UPDATE `blog-cms`.`posts` SET `title` = ?, `content` = ?, `post_type`= ? WHERE id = ?;")
	errCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(title,content,posttype,id)
	c.Redirect(http.StatusFound,"/blogar/my-post")
}

func handlePostRetrieve(c *gin.Context){
	id:=c.Param("id")
	var post api.Post


	stmt := "SELECT id, content, author FROM comments WHERE post = ?"

	rows, err := Dbase.DB.Query(stmt,id)
	if err != nil{
		_ = c.AbortWithError(500, err)
		return
	}
	defer rows.Close()

	for rows.Next(){
		var comment api.Comment
		err:= rows.Scan(&comment.ID,&comment.Content, &comment.Author)

		if err != nil{
			_ = c.AbortWithError(500, err)
			return
		}
		post.Comments = append(post.Comments, comment)
	}

	err = Dbase.DB.QueryRow("SELECT id,title,content FROM posts WHERE id=?;",id).Scan(&post.ID,&post.Title,&post.Content)
	errCheck(err)
	c.HTML(http.StatusOK, "post_detail.html", post)
}

func handlePostDelete(c *gin.Context){
	// takes the url id parameter to delete the post at the id
	id := c.Param("id")
	log.Println(id)
	del, err:= Dbase.DB.Prepare("DELETE FROM `blog-cms`.`posts` WHERE (`id`=?);")
	errCheck(err)
	defer del.Close()
	_, err = del.Exec(id)
	c.Redirect(http.StatusFound,"/blogar/my-post")
}

func handleUserPost(c *gin.Context){
	// presents the user's person post to him both private and public
	id,_:= c.Cookie("session")
	stmt := "SELECT id, title, content, post_type FROM posts WHERE author = ?"

	rows, err := Dbase.DB.Query(stmt,id)
	if err != nil{
		return
	}
	defer rows.Close()

	var posts []api.Post

	for rows.Next(){
		var p api.Post
		err:= rows.Scan(&p.ID, &p.Title,&p.Content, &p.PostType)

		if err != nil{
			return
		}
		posts = append(posts, p)

	}
	c.HTML(http.StatusOK, "userPost.html",posts)

}


// Comments Handlers

func handleCommentCreateForm(c *gin.Context){
	post_id := c.Param("id")
	content := strings.TrimSpace(c.PostForm("Content"))
	user_id,_ := c.Cookie("session")
	id := uuid.NewString()

	stmt:= "INSERT INTO `blog-cms`.`comments` (`id`,`author`, `content`, `post`) VALUES (?,?,?,?);"
	if content == "" {
		c.String(http.StatusUnauthorized,"Type in a comment before you submit")
	} else {
		prepare, err := Dbase.DB.Prepare(stmt)
		errCheck(err)
		defer prepare.Close()
		_,err = prepare.Exec(id,user_id,content,post_id)
		if err != nil {
			log.Print(err.Error())
		}


		c.Redirect(http.StatusFound,"/blogar/post/"+post_id+"/")
	}


}



