package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("C:\\Users\\Dule\\Desktop\\telep.txt") // creating...
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return
	}
	defer f.Close()
	for i := 0; i < 10; i++ { // Generating...
		_, err = f.WriteString(fmt.Sprintf("%d\r\n", i)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
}
