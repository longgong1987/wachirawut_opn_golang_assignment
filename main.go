package main

import (
	"challenge_go/internal/controllers"
	"challenge_go/internal/repositories"
	"challenge_go/internal/usecases"
	"challenge_go/services"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
)

func main() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatal("can not load env: ", err)
	}

	args := os.Args

	if len(args) <= 1 {
		log.Fatal("no Tamboon file path.")
	}

	file := args[1]

	// init opn client
	client, err := omise.NewClient(
		os.Getenv("OPN_PUBLIC_KEY"),
		os.Getenv("OPN_SECRET_KEY"),
	)
	if err != nil {
		log.Fatal("can not init opn client.")
	}
	opnClient := services.NewOpn(client)

	var (
		repositories = repositories.NewRepositories(opnClient)
		usecases     = usecases.NewUsecases(repositories)
		controllers  = controllers.NewControllers(usecases)
	)

	fmt.Printf("performing donations...\n")
	summary, err := controllers.GetDonationSummary(file)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("done.\n\n")

	fmt.Printf("%s\n", summary)

	timeUsed := time.Since(start)
	fmt.Printf("Total executed time: %d ms.\n", timeUsed.Milliseconds())

	PrintMemUsage()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc: %v MiB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc: %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys: %v MiB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC: %v\n", m.NumGC)
}
