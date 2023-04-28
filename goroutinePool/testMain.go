package main

import (
	"demo/goroutinePool"
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	ch := make(chan struct{})
	goroutinePool.Go(func() {
		fmt.Println("子协程正在执行...")
		time.Sleep(2 * time.Second) // 模拟子协程执行时间
		fmt.Println("子协程执行完成")
		ch <- struct{}{}
	})
	fmt.Println("主协程阻塞等待子协程完成...")
	<-ch
	fmt.Println("主协程继续执行...")
	cost := time.Since(start)
	fmt.Println("执行时间：", cost)
}

//func main() {
//	start := time.Now()
//
//	ch := make(chan struct{})
//	go func() {
//		fmt.Println("子协程正在执行...")
//		time.Sleep(2 * time.Second) // 模拟子协程执行时间
//		fmt.Println("子协程执行完成")
//		ch <- struct{}{}
//	}()
//	fmt.Println("主协程阻塞等待子协程完成...")
//	<-ch
//	fmt.Println("主协程继续执行...")
//	cost := time.Since(start)
//	fmt.Println("执行时间：", cost)
//}
