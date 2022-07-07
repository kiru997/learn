package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiru997/go-ex/internal/service/demo"
)

type DemoController interface {
	Finds(ctx *gin.Context)
}

type demoController struct {
	demoService demo.Service
}

func RegisterDemoController(
	r *gin.RouterGroup,
	demoService demo.Service,
) {
	g := r.Group("demo")

	var c DemoController = &demoController{
		demoService: demoService,
	}

	g.GET("finds", func(ctx *gin.Context) {
		c.Finds(ctx)
	})

}

func (d *demoController) Finds(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"number": d.demoService.Rand(c),
	})
}
