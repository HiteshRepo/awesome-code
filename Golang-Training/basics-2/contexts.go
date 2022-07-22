package main

import (
	"context"
	"fmt"
	"time"
)

func Contexts()  {
	//ctx := context.TODO()
	//doSomething(ctx)

	// developers use this to indicate that it is the starting context
	ctx := context.Background()
	doSomething(ctx)

	// Add data to context
	ctx = context.WithValue(ctx, "myKey", "myValue")
	ctx2 := context.WithValue(ctx, "myKey2", "myValue2")
	ctx3 := context.WithValue(ctx, "myKey3", "myValue3")
	ctx4 := context.WithValue(ctx, "myKey3", "myValue3")
	fmt.Println(ctx2, ctx3, ctx4)
	displayContextKeyData(ctx2, "myKey")

	updateContext(ctx, "myKey")

	// Read from context
	fmt.Printf("Contexts: myKey's value is %s\n", ctx.Value("myKey"))

	endContext()
}


func doSomething(ctx context.Context) {
	fmt.Println("Doing something!")
}

func displayContextKeyData(ctx context.Context, key string) {
	fmt.Printf("doSomething: myKey's value is %s\n", ctx.Value(key))
}

func updateContext(ctx context.Context, key string) {
	anotherContext := context.WithValue(ctx, key, "anotherValue")
	fmt.Printf("updateContext: myKey's value is %s\n", anotherContext.Value("myKey"))
}

func endContext() {
	ctx, cancelCtx := context.WithCancel(context.Background())

	printCh := make(chan int)
	go doAnother(ctx, printCh)

	for num := 1; num <= 3; num++ {
		printCh <- num
	}

	cancelCtx()

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("endContext: finished\n")
}

func doAnother(ctx context.Context, printCh <-chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("doAnother err: %s\n", err)
			}
			fmt.Printf("doAnother: finished\n")
			return
		case num := <-printCh:
			fmt.Printf("doAnother: %d\n", num)
		}
	}
}


