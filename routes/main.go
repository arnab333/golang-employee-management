package routes

import (
	"github.com/arnab333/golang-employee-management/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	// Creates a router without any middleware by default
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	// router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	api := router.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/")
	user := v1.Group("/", middlewares.VerifyToken, middlewares.CheckRole)
	role := v1.Group("/", middlewares.VerifyToken, middlewares.CheckRole)
	holiday := v1.Group("/", middlewares.VerifyToken, middlewares.CheckRole)

	authRoutes(auth)
	userRoutes(user)
	roleRoutes(role)
	holidayRoutes(holiday)

	return router
}
