package main

import (
	"context"
	"fmt"
	"github.com/n3v3rf411/go-prometheus-client-example/internal/prometheus"
	"github.com/prometheus/client_golang/api"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	_ "github.com/joho/godotenv/autoload"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.Logger())

	e.GET("/", index)
	e.GET("/stats", stats)
	e.Logger.Fatal(e.Start(":8080"))
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func stats(c echo.Context) error {
	client, err := prometheus.NewClient(api.Config{
		Address:      os.Getenv("PROMETHEUS_URL"),
		RoundTripper: api.DefaultRoundTripper,
	}, func(req *http.Request) {
		username := os.Getenv("PROMETHEUS_AUTH_USERNAME")
		password := os.Getenv("PROMETHEUS_AUTH_PASSWORD")
		if username != "" && password != "" {
			req.SetBasicAuth(username, password)
		}
	})

	if err != nil {
		return err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx,
		`up{}`,
		time.Now())
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	return c.JSON(http.StatusOK, result)
}
