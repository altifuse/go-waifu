package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Status string
type transaction struct {
	ID               string `json:"id"`
	OriginalFilename string `json:"originalFilename"`
	Status           Status `json:"status"`
}
type processEngine struct {
	Queue   []string
	Current string
	Done    []string
}

const inPath = "./in/"
const processingPath = "./processing/"
const outPath = "./out/"
const (
	Pending    Status = "PENDING"
	Processing Status = "PROCESSING"
	Done       Status = "DONE"
	Failed     Status = "FAILED"
)

var transactions = map[string]transaction{}
var runtime = processEngine{}

func main() {
	router := gin.Default()
	router.POST("/requests", uploadImage)
	router.GET("/requests/status", getAllTransactionStatus)
	router.GET("/requests/:id/status", getTransactionStatus)
	router.GET("/requests/:id/output", getOutput)
	router.Run(":8080")
}

func uploadImage(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	t := transaction{ID: uuid.New().String(), OriginalFilename: file.Filename, Status: Pending}
	log.Println("Initialized transaction " + t.ID)
	transactions[t.ID] = t
	runtime.Queue = append(runtime.Queue, t.ID)
	ctx.SaveUploadedFile(file, inPath+t.ID)
	ctx.JSON(http.StatusAccepted, t)
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
