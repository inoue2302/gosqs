package main

import (
	"./aws"
)

const CONCURRENCY = 3
const QUEUEURL = "https://sqs.ap-northeast-1.amazonaws.com/xxxxx/xxxxxx"

func main() {
	sqsr := aws.NewSqsWrapper(QUEUEURL)
	sqsr.CreateReceiveParam()
	sqsr.ReceiveMessageParallel(CONCURRENCY)
}
