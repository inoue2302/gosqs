package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"runtime"
	"sync"
	"time"
)

/**
 * AwsSqsApiラッパー
 */
type SqsWrapper struct {
	svc          *sqs.SQS
	sendParam    *sqs.SendMessageInput
	receiveParam *sqs.ReceiveMessageInput
	queueUrl     string
}

/**
 * 処理化処理
 */
func (self *SqsWrapper) initialize(queueUrl string) {
	self.queueUrl = queueUrl
	self.svc = sqs.New(&aws.Config{Region: aws.String("ap-northeast-1")})
}

/**
 * コンストラクタ
 * @return struct SqsWrapper
 */
func NewSqsWrapper(queueUrl string) *SqsWrapper {
	sqsrap := new(SqsWrapper)
	sqsrap.initialize(queueUrl)
	return sqsrap
}

/**
 * 送信用パラメター生成
 * @param msg メッセージ
 * @param queueUrl キューURL
 * @return void
 */
func (self *SqsWrapper) CreateSendParam(msg string) {
	self.sendParam = &sqs.SendMessageInput{
		MessageBody:  aws.String(msg),
		QueueUrl:     aws.String(self.queueUrl),
		DelaySeconds: aws.Int64(1),
	}
}

/**
 * メッセージ送信
 * @return void
 */
func (self *SqsWrapper) SendMesage() {
	resp, err := self.svc.SendMessage(self.sendParam)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}

/**
 * 受信パラメーター生成
 * @return void
 */
func (self *SqsWrapper) CreateReceiveParam() {
	self.receiveParam = &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(self.queueUrl),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}
}

/**
 * メッセージ受信
 * @return void
 */
func (self *SqsWrapper) ReceiveMessage() {
	resp, err := self.svc.ReceiveMessage(self.receiveParam)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if resp == nil {
		fmt.Println("emptye message")
		return
	}
	if len(resp.Messages) == 0 {
		fmt.Println("queue is empty")
		return
	}
	message := resp.Messages[0].Body
	handle := resp.Messages[0].ReceiptHandle
	fmt.Println(*message)
	self.DeleteMessage(*handle)
	time.Sleep(time.Second)
}

/**
 * メッセージ並列受信
 * @param int concurrency
 * @return void
 */
func (self *SqsWrapper) ReceiveMessageParallel(concurrency int) {
	cpus := runtime.NumCPU()
	semaphore := make(chan int, cpus)
	var w sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			semaphore <- 1
			self.ReceiveMessage()
			<-semaphore
		}()
	}
	w.Wait()
}

/**
 * メッセージ削除
 * @param string handle
 * @return void
 */
func (self *SqsWrapper) DeleteMessage(handle string) {
	param := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(self.queueUrl), // Required
		ReceiptHandle: aws.String(handle),        // Required
	}
	self.svc.DeleteMessage(param)
}
