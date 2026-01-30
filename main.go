package main

import (
	"fmt"
	"log"
	"messenger/handlers"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	sqlDB, _ := db.DB()
	if err := goose.Up(sqlDB, "./migrations"); err != nil {
		log.Fatalf("goose migration failed: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/chats/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/chats/" && r.Method == http.MethodPost {
			handlers.CreateChatHandler(db)(w, r)
			return
		}

		if (strings.HasSuffix(path, "/messages/") || strings.HasSuffix(path, "/messages")) && r.Method == http.MethodPost {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}

			r2 := *r
			urlCopy := *r.URL
			urlCopy.Path = path
			r2.URL = &urlCopy

			handlers.SendMessageHandler(db)(w, &r2)
			return
		}

		if r.Method == http.MethodGet && len(path) > len("/chats/") {
			handlers.GetChatHandler(db)(w, r)
			return
		}

		if r.Method == http.MethodDelete && len(path) > len("/chats/") {
			handlers.DeleteChatHandler(db)(w, r)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	mux.Handle("/docs", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/openapi.json",
		Path:    "/docs",
		Title:   "Chat API",
	}, nil))

	mux.Handle("/docs/", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/openapi.json",
		Path:    "/docs/",
		Title:   "Chat API",
	}, nil))

	mux.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.json")
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
