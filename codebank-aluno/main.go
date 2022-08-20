package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/codeedu/codebank/infrastructure/grpc/server"
	"github.com/codeedu/codebank/infrastructure/kafka"
	"github.com/codeedu/codebank/infrastructure/repository"
	"github.com/codeedu/codebank/usecase"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello Go")
	db := setupDb()
	defer db.Close()
	producer := setupKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	serveGrpc(processTransactionUseCase)

}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer
	return useCase
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer("host.docker.internal:9094")
	return producer
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db", "5432", "postgres", "root", "codebank")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}

func serveGrpc(processTransactionUsecase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUsecase
	fmt.Println("Rodando gRPC Server")
	grpcServer.Serve()
}

// cc := domain.NewCreditCard()
// cc.Number = "5594 8871 6180 2534"
// cc.Name = "Tiago Wesley"
// cc.ExpirationYear = 2024
// cc.ExpirationMonth = 3
// cc.CVV = 569
// cc.Limit = 1000
// cc.Balance = 0

// repo := repository.NewTransactionRepositoryDb(db)
// err := repo.CreateCreditCard(*cc)
// if err != nil {
// 	fmt.Println(err)
// }
