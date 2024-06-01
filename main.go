package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Status string
type Transaction struct {
	ID               string `json:"id"`
	OriginalFilename string `json:"originalFilename"`
	Status           Status `json:"status"`
}
type Response struct {
	Message string `json:"message"`
}

const inPath = "./in/"
const processingPath = "./processing/"
const outPath = "./out/"
const failedPath = "./failed/"
const (
	Pending    Status = "PENDING"
	Processing Status = "PROCESSING"
	Done       Status = "DONE"
	Failed     Status = "FAILED"
)
const transactionQueueCapacity = 10

var transactions = map[string]Transaction{}
var transactionQueue chan string

func main() {
	transactionQueue = make(chan string, transactionQueueCapacity)
	go runtimeLoop(transactionQueue)
	router := gin.Default()
	router.POST("/requests", uploadImage)
	router.GET("/requests/status", getAllTransactionStatus)
	router.GET("/requests/:id/status", getTransactionStatus)
	router.GET("/requests/:id/output", getOutput)
	router.Run(":8080")
}

func uploadImage(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	t := Transaction{ID: uuid.New().String(), OriginalFilename: file.Filename, Status: Pending}
	log.Println("Preparing transaction " + t.ID)
	if enqueueTransaction(t) {
		transactions[t.ID] = t
		log.Println("Transaction enqueued successfully")
		ctx.SaveUploadedFile(file, inPath+t.ID)
		ctx.JSON(http.StatusAccepted, t)
	} else {
		log.Println("Transaction could not be enqueued - queue is full")
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, Response{Message: "Queue is currently full, please try again later."})
	}
}

func getAllTransactionStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, transactions)
}

func getTransactionStatus(ctx *gin.Context) {
	t, ok := transactions[ctx.Param("id")]
	if ok {
		ctx.JSON(http.StatusOK, t)
	} else {
		ctx.Status(http.StatusNotFound)
	}
}

func getOutput(ctx *gin.Context) {
	t, ok := transactions[ctx.Param("id")]
	if !ok || (ok && t.Status != Done) {
		ctx.Status(http.StatusNotFound)
	}
	ctx.File(outPath + t.ID)
}

func enqueueTransaction(transaction Transaction) bool {
	select {
	case transactionQueue <- transaction.ID:
		return true
	default:
		return false
	}
}

func runtimeLoop(jobChannel <-chan string) {
	for {
		log.Println("Looping over transaction queue")
		for transactionID := range jobChannel {
			processTransaction(transactionID)
		}
	}
}

func processTransaction(transactionID string) {
	log.Println("Processing transaction " + transactionID)
	// TODO: move file to processing dir
	// TODO: call waifu2x
	// TODO: if exit code 0, confirm file is present in out dir
	// - if yes, delete file in processing and set status to Done;
	//   if no, move file to failed dir and set status to Failed
	time.Sleep(5 * time.Second)
	t := transactions[transactionID]
	t.Status = Done
	transactions[transactionID] = t
}
