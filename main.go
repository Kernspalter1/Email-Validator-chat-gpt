package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("=== LocalEmailHealthChecker START ===")

	fmt.Println("Args:", os.Args)

	// TODO: eigentliche Logik kommt hier rein
	fmt.Println("Program reached end of main()")

	fmt.Println("Waiting 10 seconds before exit...")
	time.Sleep(10 * time.Second)
}
