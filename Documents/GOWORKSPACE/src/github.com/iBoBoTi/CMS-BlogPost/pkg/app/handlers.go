package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iBoBoTi/CMS-BlogPost/pkg/api"
	"log"
	"net/http"
	"strings"
	//"github.com/iBoBoTi/CMS-BlogPost/pkg/api"
)

func (s *Server) index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (s *Server) home(c *gin.Context) {
	stmt := `SELECT posts.id, posts.title, posts.content,posts.post_type, users.user_name FROM posts INNER JOIN users ON posts.author = users.id WHERE posts.post_type = 'Public'`

	rows, err := s.DB.DB.Query(stmt)
	if err != nil {
		return
	}
	defer rows.Close()

	var posts []api.Post

	for rows.Next() {
		var p api.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.PostType, &p.UserName)

		if err != nil {
			return
		}
		posts = append(posts, p)

	}
	c.HTML(http.StatusOK, "home.html", posts)
}

func (s *Server) handleSignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "SignUp.html", nil)
}

func (s *Server) handleSignUpAuth(c *gin.Context) {
	// function will handle the form submitted at signup and add the details to the database
	// validate form takes in a form type return bool true results in processing the details while false results
	// in redirecting to the error string
	checkString, check := s.ValidateSignUpForm(c)

	if check == true {
		// process form and add to database
		id := uuid.NewString()
		firstname := c.PostForm("FirstName")
		lastname := c.PostForm("LastName")
		username := strings.Trim(c.PostForm("UserName"), " ")
		email := c.PostForm("Email")
		password, _ := HashPassword(c.PostForm("password"))

		//insertDynStmt := `insert into "Students"("Name", "Roll") values($1, $2)`
		//_, e = db.Exec(insertDynStmt, "Jane", 2)
		//CheckError(e)

		stmt := `insert into "users" ("id", "first_name", "last_name", "user_name", "email", "password") values ($1,$2,$3,$4,$5,$6)`
		prepare, err := s.DB.DB.Prepare(stmt)
		//errCheck(err)
		if err != nil {
			panic(err)
		}
		defer prepare.Close()
		_, err = prepare.Exec(id, firstname, lastname, username, email, password)
		if err != nil {
			log.Print(err.Error())
		}
		// redirects to login
		c.Redirect(http.StatusFound, "/login")

	} else {
		c.String(http.StatusBadRequest, checkString)
	}

}

func (s *Server) handleLogin(c *gin.Context) {
	// function will present the login form to the user
	c.HTML(http.StatusOK, "login.html", nil)
}

func (s *Server) handleLoginAuth(c *gin.Context) {
	// function will collect details from the login form
	// function will verify details as well as create a session if verification passes
	loginCheck, sessionToken := s.ValidateLoginForm(c)

	if loginCheck == true {
		c.SetCookie("session", sessionToken, 3600, "/", "https://blogaar.herokuapp.com", false, true)
		c.Redirect(http.StatusFound, "/blogar")
	} else {
		c.String(http.StatusUnauthorized, "Details doesn't exist")
	}
}

func (s *Server) handleLogOut(c *gin.Context) {
	// function will delete our session and redirect to the home page
	c.SetCookie("session", "", -1, "/", "https://blogaar.herokuapp.com", true, true)
	c.Redirect(http.StatusFound, "/")
}

func (s *Server) handlePostCreate(c *gin.Context) {
	// presents the form for creating a post
	c.HTML(http.StatusOK, "CreateBlogPostForm.html", nil)
}

func (s *Server) handlePostCreateForm(c *gin.Context) {
	// processes the creat blog post form and adds it to the data
	id := uuid.NewString()
	title := c.PostForm("Title")
	content := c.PostForm("Content")
	posttype := c.PostForm("PostType")
	username, _ := c.Cookie("session")

	stmt := fmt.Sprint("INSERT INTO posts (id,author,title, content, post_type) VALUES ($1,$2,$3,$4,$5);")
	if title == "" || content == "" {
		c.String(http.StatusUnauthorized, "Please check your entry dear Blogar")
	} else {
		prepare, err := s.DB.DB.Prepare(stmt)
		errCheck(err)
		defer prepare.Close()
		_, err = prepare.Exec(id, username, title, content, posttype)
		if err != nil {
			log.Print(err.Error())
		}
		c.Redirect(http.StatusFound, "/blogar/")
	}

}

func (s *Server) handlePostEdit(c *gin.Context) {
	id := c.Param("id")

	row := s.DB.DB.QueryRow("SELECT title,content,post_type,id FROM posts WHERE id=$1;", id)
	var post api.Post
	err := row.Scan(&post.Title, &post.Content, &post.PostType, &post.ID)
	errCheck(err)

	c.HTML(http.StatusOK, "EditBlogPostForm.html", post)
}

func (s *Server) handlePostEditForm(c *gin.Context) {
	id := c.Param("id")

	title := c.PostForm("Title")
	content := c.PostForm("Content")
	posttype := c.PostForm("PostType")

	stmt, err := s.DB.DB.Prepare("UPDATE posts SET title = $1, content = $2, post_type= $3 WHERE id = $4;")
	errCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(title, content, posttype, id)
	c.Redirect(http.StatusFound, "/blogar/my-post")
}

func (s *Server) handlePostRetrieve(c *gin.Context) {
	id := c.Param("id")
	var post api.Post

	stmt := "SELECT c.id, c.content, u.user_name as author FROM comments c JOIN users u ON c.author = u.id AND c.post_id = $1"

	rows, err := s.DB.DB.Query(stmt, id)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment api.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.Author)

		if err != nil {
			_ = c.AbortWithError(500, err)
			return
		}
		post.Comments = append(post.Comments, comment)
	}

	err = s.DB.DB.QueryRow("SELECT id,title,content FROM posts WHERE id=$1", id).Scan(&post.ID, &post.Title, &post.Content)
	fmt.Println(post)
	errCheck(err)
	c.HTML(http.StatusOK, "post_detail.html", post)
}

func (s *Server) handlePostDelete(c *gin.Context) {
	// takes the url id parameter to delete the post at the id
	id := c.Param("id")
	log.Println(id)
	del, err := s.DB.DB.Prepare("DELETE FROM posts WHERE (id=$1);")
	errCheck(err)
	defer del.Close()
	_, err = del.Exec(id)
	c.Redirect(http.StatusFound, "/blogar/my-post")
}

func (s *Server) handleUserPost(c *gin.Context) {
	// presents the user's person post to him both private and public
	id, _ := c.Cookie("session")
	stmt := "SELECT id, title, content, post_type FROM posts WHERE author = $1"

	rows, err := s.DB.DB.Query(stmt, id)
	if err != nil {
		return
	}
	defer rows.Close()

	var posts []api.Post

	for rows.Next() {
		var p api.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.PostType)

		if err != nil {
			return
		}
		posts = append(posts, p)

	}
	c.HTML(http.StatusOK, "userPost.html", posts)

}

// Comments Handlers

func (s *Server) handleCommentCreateForm(c *gin.Context) {
	post_id := c.Param("id")
	content := strings.TrimSpace(c.PostForm("Content"))
	user_id, _ := c.Cookie("session")
	id := uuid.NewString()

	stmt := "INSERT INTO comments (id,author, content, post_id) VALUES ($1,$2,$3,$4);"

	if content == "" {
		c.String(http.StatusUnauthorized, "Type in a comment before you submit")
	} else {
		prepare, err := s.DB.DB.Prepare(stmt)
		errCheck(err)
		defer prepare.Close()

		_, err = prepare.Exec(id, user_id, content, post_id)
		if err != nil {
			log.Print(err.Error())
		}

		c.Redirect(http.StatusFound, "/blogar/post/"+post_id+"/")
	}

}
