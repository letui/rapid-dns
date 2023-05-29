package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"rapid-dns/db"
	"rapid-dns/model"
	"time"
)

func Login(c *gin.Context) {
	if HasSession(c) {
		c.JSON(http.StatusOK, OK(time.Now().Unix()))
		return
	}
	user := db.User{}
	if c.ShouldBind(&user) != nil {
		c.JSON(http.StatusOK, Error("非法的请求数据"))
		return
	}
	if db.ValidateUser(user) {
		session := sessions.Default(c)
		session.Set("USER", user.Name)
		SaveSession(c, user.Name)
		c.JSON(http.StatusOK, OK(time.Now().Unix()))
	} else {
		c.JSON(http.StatusOK, Error("账号密码不匹配"))
	}
}
func Logout(c *gin.Context) {
	if HasSession(c) {
		ClearSession(c)
	}
	c.JSON(http.StatusOK, OK(time.Now().Unix()))
}

func Register(c *gin.Context) {
	user := db.User{}
	if c.ShouldBind(&user) != nil {
		c.JSON(http.StatusOK, Error("请求格式有误"))
		return
	}
	if user.Password != user.RetypePassword {
		c.JSON(http.StatusOK, Error("两次密码输入不一致"))
		return
	}
	if db.ExistSameUsername(user) {
		c.JSON(http.StatusOK, Error("存在相同用户名"))
		return
	}
	db.AddNewUser(user)
	c.JSON(http.StatusOK, OK(time.Now().Unix()))
}

func UpdatePassword(c *gin.Context) {
	if !HasSession(c) {
		c.JSON(http.StatusOK, Error("请登录后操作"))
		return
	}
	session := sessions.Default(c)
	username := session.Get("USER").(string)
	password := model.PasswordRequest{}
	if c.ShouldBind(&password) != nil {
		c.JSON(http.StatusOK, Error("请求格式有误"))
		return
	}
	if password.NewPassword != password.NewRPassword {
		c.JSON(http.StatusOK, Error("两次密码输入不一致"))
		return
	}
	user := db.User{}
	user.Name = username
	user.Password = password.OldPassword
	if db.ValidateUser(user) {
		user.Password = password.NewRPassword
		db.AddNewUser(user)
		c.JSON(http.StatusOK, OK(time.Now().Unix()))
	} else {
		c.JSON(http.StatusOK, Error("原密码输入有误"))
	}
}
