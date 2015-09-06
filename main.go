package main

import (
	"./aws"
)

const QUEUEURL = "https://sqs.ap-northeast-1.amazonaws.com/124933150345/sqstest"

func send_message(message string) int {
	sqsr := aws.NewSqsWrapper(QUEUEURL)
	sqsr.CreateSendParam(message)
	sqsr.SendMesage()
	return 1
}

func receive() {
	sqsr := aws.NewSqsWrapper(QUEUEURL)
	sqsr.CreateReceiveParam()
	sqsr.ReceiveMessage()
}

func main() {
	send := make(chan int)
	quit := make(chan bool)
	workerquit := make(chan bool)

	go func() {
	loop:
		for {
			select {
			case <-quit:
				workerquit <- true
				break loop
			case <-send:
				receive()
			}
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			send <- send_message("HelloSqs")
		}
		quit <- true
	}()
	<-workerquit
}
