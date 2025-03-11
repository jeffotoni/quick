package main

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MainController struct {
	web.Controller
}

func (c MainController) Post() {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Erro ao decodificar JSON"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}

func main() {
	web.Router("/v1/user", &MainController{}, "post:Post")
	web.Run()
}
