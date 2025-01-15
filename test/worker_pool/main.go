package main

import (
	"fmt"
	"time"

	"github.com/lastnight2022-cc/tools_lib/utils/worker_pool"
)

func main() {
	workerCount := 5
	wp, err := worker_pool.NewWorkerPool(workerCount)
	if err != nil {
		fmt.Println("Failed to create worker pool:", err)
		return
	}
	defer wp.Close()

	// 提交一些任务到 WorkerPool 中。
	for i := 0; i < 10; i++ {
		task := func(i int) worker_pool.Task {
			return func() error {
				fmt.Printf("Executing task %d\n", i)
				time.Sleep(time.Second) // 模拟耗时操作
				return nil
			}
		}(i)

		if err := wp.Submit(task); err != nil {
			fmt.Println("Failed to submit task:", err)
		}
	}

	// 等待所有任务完成（可选）
	time.Sleep(2 * time.Second) // 这里仅用于演示目的；实际应用中可以有其他方式确保任务完成

	fmt.Println("All tasks submitted. Closing the worker pool.")
}
