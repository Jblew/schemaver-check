package main

import (
	"fmt"
	"log"
	"time"
)

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		fmt.Println("This is attempt number", i)
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		err = f()
		if err != nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
