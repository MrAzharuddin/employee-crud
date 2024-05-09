package main

import (
	"context"
	"github.com/MrAzharuddin/employee-crud/employee-service/config"
	_ "github.com/MrAzharuddin/employee-crud/employee-service/docs"
	restcontrollers "github.com/MrAzharuddin/employee-crud/employee-service/pkg/rest/server/controllers"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinhashubham95/go-actuator"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"os"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func ServeRoutes() *gin.Engine {
	router := gin.Default()
	employeeController, err := restcontrollers.NewEmployeeController()
	if err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := router.Group("/v1")
	{

		v1.POST("/employees", employeeController.CreateEmployee)

		v1.GET("/employees/:id", employeeController.FetchEmployee)

		v1.GET("/employees", employeeController.FetchEmployees)

		v1.PUT("/employees/:id", employeeController.UpdateEmployee)

		v1.DELETE("/employees/:id", employeeController.DeleteEmployee)

		v1.POST("/employees/random", employeeController.PushEmployee)

	}
	return router
}

//	@title			employee-service
//	@version		1.0
//	@description	Testing Swagger APIs.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8000
//	@BasePath		/v1
//	@schemes		http
func main() {

	// rest server configuration
	router := ServeRoutes()
	var restTraceProvider *sdktrace.TracerProvider
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// add opentel
		restTraceProvider = config.InitRestTracer(serviceName, collectorURL, insecure)
		router.Use(otelgin.Middleware(serviceName))
	}
	defer func() {
		if restTraceProvider != nil {
			if err := restTraceProvider.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()
	// add actuator
	addActuator(router)
	// add prometheus
	addPrometheus(router)

	Port := ":8000"
	log.Println("Server started")
	if err := router.Run(Port); err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func addPrometheus(router *gin.Engine) {
	router.GET("/metrics", prometheusHandler())
}

func addActuator(router *gin.Engine) {
	actuatorHandler := actuator.GetActuatorHandler(&actuator.Config{Endpoints: []int{
		actuator.Env,
		actuator.Info,
		actuator.Metrics,
		actuator.Ping,
		// actuator.Shutdown,
		actuator.ThreadDump,
	},
		Env:     "dev",
		Name:    "employee-service",
		Port:    8000,
		Version: "0.0.1",
	})
	ginActuatorHandler := func(ctx *gin.Context) {
		actuatorHandler(ctx.Writer, ctx.Request)
	}
	router.GET("/actuator/*endpoint", ginActuatorHandler)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
