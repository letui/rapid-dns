package web

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"time"
)

const SESSION_KEY = "X-SESSION"
const SESSION_ID = "X-SESSION-ID"

var Port = "8053"

type Server struct {
}

type Resp struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

//go:embed ui/index.html ui/css/* ui/js/* ui/fonts/* ui/*
var vuefs embed.FS

func (srv *Server) StartWeb() {
	gin.SetMode(gin.ReleaseMode)
	grouter := gin.New()
	must := template.Must(template.New("").ParseFS(vuefs, "ui/*.html"))
	grouter.SetHTMLTemplate(must)
	grouter.GET("/", func(context *gin.Context) {
		context.Redirect(302, "/ui")
	})

	grouter.Use(EnableCookieSession()) // cookie as store
	grouter.GET("/logo.png", func(c *gin.Context) {
		c.FileFromFS("ui/logo.png", http.FS(vuefs))
	})
	base, _ := fs.Sub(vuefs, "ui")
	grouter.StaticFS("ui", http.FS(base))
	js, _ := fs.Sub(vuefs, "ui/js")
	grouter.StaticFS("js", http.FS(js))
	css, _ := fs.Sub(vuefs, "ui/css")
	grouter.StaticFS("css", http.FS(css))
	fonts, _ := fs.Sub(vuefs, "ui/fonts")
	grouter.StaticFS("fonts", http.FS(fonts))

	api := grouter.Group("/api")
	api.POST("/auth/login", Login)
	api.POST("/auth/logout", Logout)
	api.POST("/auth/updatePassword", UpdatePassword)
	api.POST("/auth/register", Register)

	api.GET("/domain/supported", Supported)
	api.GET("/domain/search", QueryForRegister)
	api.POST("/domain/new", AddDomain)
	api.POST("/domain/drop", DeleteDomain)
	api.POST("/domain/modify", ModifyDomain)
	api.POST("/domain/list", ListDomains)

	grouter.Run("0.0.0.0:" + Port)
}

func EnableCookieSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(SESSION_KEY))
	return sessions.Sessions(SESSION_KEY, store)
}

// register and login will save session
func SaveSession(ctx *gin.Context, username string) string {
	session := sessions.Default(ctx)
	// simple encrypt: sha256(usename + string(timestamp))
	contact := username + strconv.FormatInt(time.Now().UnixMilli(), 10)
	sum := sha256.Sum256([]byte(contact))
	val := fmt.Sprintf("%x", sum)
	session.Set(SESSION_ID, val)
	session.Save()
	sessionID := session.Get(SESSION_ID).(string)
	ctx.Writer.Header().Set(SESSION_ID, sessionID)
	return sessionID
}

// logout will clear session
func ClearSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}

func HasSession(ctx *gin.Context) bool {
	session := sessions.Default(ctx)
	if val := session.Get(SESSION_ID); val == nil {
		return false
	}
	return true
}

func GetSession(ctx *gin.Context) string {
	session := sessions.Default(ctx)
	val := session.Get(SESSION_ID)
	if val.(string) == "" {
		return ""
	}
	return val.(string)
}

func OK(obj any) Resp {
	return Resp{200, "ok", obj}
}

func Error(obj any) Resp {
	return Resp{500, "error", obj}
}
