package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// BestPractices and Caveats
// Following is a list of best practices that you can follow while using a context.
// - Do not store a context within a struct type
// - Context should flow through your program.
// 		For example, in case of an HTTP request, a new context can be created for each incoming request
//		which can be used to hold a request_id or put some common information in the context like currently
//		logged in user which might be useful for that particular request.
// - Always pass context as the first argument to a function.
// - Whenever you are not sure whether to use the context or not, it is better to use the context.ToDo() as a placeholder.
// - Only the parent goroutine or function should the cancel context.
//		Therefore do not pass the cancelFunc to downstream goroutines or functions.
//		Golang will allow you to pass the cancelFunc around to child goroutines but it is not a recommended practice.

func main() {
	// Passing value in the context
	// curl -v http://localhost:8080/welcome
	// helloWorldHandler := http.HandlerFunc(HelloWorld)
	// http.Handle("/welcome", injectMsgID(helloWorldHandler))
	// http.ListenAndServe(":8080", nil)

	// Playing with the cancel function
	// ctx := context.Background()
	// cancelCtx, cancelFunc := context.WithCancel(ctx)
	// go task(cancelCtx)
	// time.Sleep(time.Second * 3)
	// cancelFunc()
	// time.Sleep(time.Second * 1)

	// // Playing with timeout function
	// ctx := context.Background()
	// cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	// defer cancel()
	// go task(cancelCtx)
	// time.Sleep(time.Second * 4)

	// Playing with deadline function
	ctx := context.Background()
	cancelCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()
	go task(cancelCtx)
	time.Sleep(time.Second * 6)
}

//HelloWorld hellow world handler
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	msgID := ""
	if m := r.Context().Value("msgId"); m != nil {
		if value, ok := m.(string); ok {
			msgID = value
		}
	}
	w.Header().Add("msgId", msgID)
	w.Write([]byte("Hello, world"))
}

// Inject msg id in the request context
func injectMsgID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "msgId", msgID)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func task(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully exit")
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
