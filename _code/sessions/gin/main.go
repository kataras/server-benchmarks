package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"

	uuid "github.com/iris-contrib/go.uuid"
)

func main() {
	gin.SetMode("release")

	app := gin.New()
	storeKey := securecookie.GenerateRandomKey(32)
	store := memstore.NewStore(storeKey)
	app.Use(sessions.Sessions("session", store))

	app.GET("/sessions", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		id, _ := uuid.NewV4()
		session.Set("ID", id.String())

		session.Set("name", "John Doe")
		if err := session.Save(); err != nil {
			ctx.String(500, err.Error())
			return
		}

		v := session.Get("name")
		name, ok := v.(string)
		if !ok {
			ctx.String(500, "server error")
			return
		}

		ctx.String(200, name)
	})

	app.Run(":5000")
}
