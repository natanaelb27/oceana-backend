package routers

import (
	"fmt"
	"log"
	"os"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(h *handlers.Handler) {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	r.Use(cors.New(config))

	r.GET("/", h.Root)
	r.Static("/docs", "cmd/docs")

	public := r.Group("/")
	private := r.Group("/", middlewares.Auth)

	role(public, private, h)
	profile(public, private, h)
	auth(public, private, h)
	category(public, private, h)
	address(private, h)
	shop(public, private, h)
	myShop(public, private, h)
	cart(public, private, h)
	product(public, private, h)
	review(public, private, h)
	subDistrict(public, private, h)
	rajaOngkir(public, private, h)
	marketplaceVoucher(public, private, h)
	courier(public, private, h)
	city(public, private, h)

	r.NoRoute(h.NotFound)

	err := r.Run(fmt.Sprintf(":%v", os.Getenv("ROUTER_PORT")))
	if err != nil {
		log.Fatal("Routers initiation failed")
	}
}
