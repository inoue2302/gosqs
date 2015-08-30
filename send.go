package main

import (
	"./aws"
	"flag"
)

const QUEUEURL = "https://sqs.ap-northeast-1.amazonaws.com/xxxxxx/xxxxxx"

func main() {
	var message = flag.String("msg", "HelloSqs", "help message for SqsSendMessage")
	flag.Parse()
	sqsr := aws.NewSqsWrapper(QUEUEURL)
	sqsr.CreateSendParam(*message)
	sqsr.SendMesage()
}
