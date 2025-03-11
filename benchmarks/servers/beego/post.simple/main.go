package main

import (
	"fmt"
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

func (c *MainController) Post() {
	var user User
	if err := c.BindJSON(&user); err != nil {
		fmt.Println("Erro ao decodificar JSON:", err) // Log do erro
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Erro ao decodificar JSON: " + err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}

func main() {
	// Set the port to 8081
	//web.BConfig.Listen.HTTPPort = 8081

	// Set the run mode manually (avoids loading the app.conf file)
	web.BConfig.RunMode = "prod" // or "dev"

	web.Router("/v1/user", &MainController{}, "post:Post")
	web.Run()
}
