package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Korpenter/club/internal/app"
	"github.com/Korpenter/club/internal/config"
	"github.com/Korpenter/club/internal/handler"
	"github.com/Korpenter/club/internal/service"
	"github.com/Korpenter/club/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path_to_input_file>", os.Args[0])
	}

	inputFilePath := os.Args[1]
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", inputFilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cfg, err := config.NewConfig(scanner)
	if err != nil {
		fmt.Println(err)
		return
	}
	repo := storage.NewInMemRepo(cfg)
	service := service.New(cfg, repo)
	handler := handler.NewFileHandler(scanner, service, cfg)
	app := app.NewApp(handler)
	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
