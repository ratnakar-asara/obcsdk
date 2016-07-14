package main

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

func main() {
	t := time.Now()
	//fmt.Println(t)
	//fmt.Println(t.Format(time.RFC3339))
	var buf bytes.Buffer
	//logger := log.New(&buf, t.Format(time.RFC3339) +" ", log.Lshortfile)
	logger := log.New(&buf, log.Ldate+log.Ltime, log.Lshortfile)
	logger.Print("Hello, log file!")

	fmt.Print(&buf)
}
