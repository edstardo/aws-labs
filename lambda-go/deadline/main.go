package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Deadline  int `json:"deadline"`
	Execution int `json:"execution"`
}

func HandleRequest(ctx context.Context, req Request) error {

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	deadline := make(chan interface{})
	go func() {
		<-time.After(time.Duration(req.Deadline) * time.Second)
		close(deadline)
	}()

	completion := make(chan interface{})
	go func() {
		<-time.After(time.Duration(req.Execution) * time.Second)
		close(completion)
	}()

	for {
		select {
		case <-deadline:
			return errors.New("Deadline reached before handler completed!")
		case <-completion:
			fmt.Println("Handler completed execution")
			return nil
		case <-ctx.Done():
			fmt.Println("Context channel closed!")
			return errors.New("Context channel closed before handler completion!")
		case <-ticker.C:
			fmt.Println("Handler running...")
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
