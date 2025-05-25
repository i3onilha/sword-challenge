package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"sword-challenge/config"
	"sword-challenge/internal/controllers"
	"sword-challenge/internal/middleware"
	"sword-challenge/internal/repository/mysql"
	"sword-challenge/internal/service"
	"sword-challenge/pkg/messaging"

	docs "sword-challenge/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

// @title           Sword Challenge API
// @version         1.0
// @description     A task management API for technicians and managers.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// --- Provide Functions ---

func newRouter() *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
            <!DOCTYPE html>
            <html lang="en">
            <head>
              <meta charset="UTF-8" />
              <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
              <title>Welcome to Sword Challenge API</title>
              <style>
                body {
                  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                  background-color: #111;
                  color: #fff;
                  text-align: center;
                  padding: 100px 20px;
                }

                h1 {
                  font-size: 3em;
                  margin-bottom: 10px;
                  color: #39ff14;
                }

                p {
                  font-size: 1.2em;
                  margin-bottom: 40px;
                  color: #ccc;
                }

                .btn {
                  background-color: #39ff14;
                  color: #000;
                  padding: 15px 30px;
                  font-size: 1em;
                  border: none;
                  border-radius: 5px;
                  text-decoration: none;
                  cursor: pointer;
                  transition: background 0.3s ease;
                }

                .btn:hover {
                  background-color: #2ecc71;
                }
              </style>
            </head>
            <body>
              <h1>Welcome to Sword Challenge API</h1>
              <p>Built for speed. Scalable and efficient task management.</p>
              <a href="/swagger/index.html" class="btn">View API Docs (Swagger)</a>
            </body>
            </html>
        `))
	})

	return router
}

func newMessageBroker() (messaging.MessageBroker, error) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable is not set")
	}
	return messaging.NewRabbitMQ(rabbitmqURL)
}

// --- Run server lifecycle ---

func runServer(lc fx.Lifecycle, router *gin.Engine, consumer *messaging.NotificationConsumer) {
	port := config.GetEnv("PORT", "3000")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start notification consumer
			if err := consumer.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := router.Run(":" + port); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return nil
		},
	})
}

// --- Route Registration ---

func registerRoutes(
	router *gin.Engine,
	taskController *controllers.TaskController,
	notificationController *controllers.NotificationController,
) {
	// Create middleware instances
	authMiddleware := middleware.GinAuthMiddleware()

	tasks := router.Group("/api/tasks")
	tasks.Use(authMiddleware)
	{
		// Technician routes - can only access their own tasks
		tasks.POST("", middleware.RequireRole("technician"), taskController.CreateTask)
		tasks.GET("", middleware.RequireRole("technician", "manager"), taskController.GetTasks)    // Both roles can access, but service layer filters results
		tasks.GET("/:id", middleware.RequireRole("technician", "manager"), taskController.GetTask) // Both roles can access, but service layer filters results
		tasks.PUT("/:id", middleware.RequireRole("technician"), taskController.UpdateTask)
		tasks.DELETE("/:id", middleware.RequireRole("manager"), taskController.DeleteTask)
	}

	notifications := router.Group("/api/notifications")
	notifications.Use(authMiddleware)
	{
		notifications.GET("", middleware.RequireRole("manager"), notificationController.GetUnreadNotifications)
		notifications.PUT("/:id/read", middleware.RequireRole("manager"), notificationController.MarkAsRead)
	}
}

// --- Main Function ---

func main() {
	app := fx.New(
		// Providers
		fx.Provide(
			config.InitDB,
			mysql.NewUserRepository,
			mysql.NewTaskRepository,
			mysql.NewNotificationRepository,
			newMessageBroker,
			service.NewTaskService,
			service.NewNotificationService,
			controllers.NewTaskController,
			controllers.NewNotificationController,
			newRouter,
			messaging.NewNotificationConsumer,
		),
		// Invokes
		fx.Invoke(registerRoutes, runServer),
	)

	app.Run()
}
