package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"rapid-dns/db"
	"rapid-dns/model"
)

func ListDomains(c *gin.Context) {
	if !HasSession(c) {
		c.JSON(200, Error("请登录后操作"))
		return
	}
	session := sessions.Default(c)
	username := session.Get("USER").(string)
	domains := db.DomainListOfUser(username)
	var list []model.DomainRequest
	for _, o := range domains {
		item := model.DomainRequest{Name: o, Ipv4: string(db.Query(o + "."))}
		list = append(list, item)
	}
	c.JSON(200, OK(list))
}

func ModifyDomain(c *gin.Context) {
	if !HasSession(c) {
		c.JSON(200, Error("请登录后操作"))
		return
	}
	request := new(model.DomainRequest)
	if c.ShouldBind(&request) != nil {
		c.JSON(200, Error("请求格式不正确"))
		return
	}
	session := sessions.Default(c)
	username := session.Get("USER").(string)
	if db.OwnDomain(request.Name, username) {
		added := db.AddDomainWithIpv4(request.Name, request.Ipv4)
		if added {
			c.JSON(200, OK("更新域名成功"))
		}
	} else {
		c.JSON(200, Error("更新域名失败"))
	}
}

func AddDomain(c *gin.Context) {
	if !HasSession(c) {
		c.JSON(200, Error("请登录后操作"))
		return
	}
	request := new(model.DomainRequest)
	if c.ShouldBind(&request) != nil {
		c.JSON(200, Error("请求格式不正确"))
		return
	}
	if db.ExistDomain(request.Name) {
		c.JSON(200, Error("域名已被注册，请换一个进行尝试"))
		return
	}
	added := db.AddDomainWithIpv4(request.Name, request.Ipv4)
	if added {
		session := sessions.Default(c)
		username := session.Get("USER").(string)
		db.MarkDomainForUser(request.Name, username)
		c.JSON(200, OK("添加域名成功"))
	} else {
		c.JSON(200, Error("添加域名失败"))
	}
}

func DeleteDomain(c *gin.Context) {
	if !HasSession(c) {
		c.JSON(200, Error("请登录后操作"))
		return
	}
	request := new(model.DomainRequest)
	if c.ShouldBind(&request) != nil {
		c.JSON(200, Error("请求格式不正确"))
		return
	}
	deleted := db.DeleteDomainWithIpv4(request.Name)
	session := sessions.Default(c)
	username := session.Get("USER").(string)
	if deleted {
		db.UnMarkDomainForUser(request.Name, username)
		c.JSON(200, OK("删除域名成功"))
	} else {
		c.JSON(200, Error("删除域名失败"))
	}
}

func QueryForRegister(c *gin.Context) {
	value := c.Query("name")
	domains := db.ListUsableDomains(value)
	c.JSON(200, OK(domains))
}
