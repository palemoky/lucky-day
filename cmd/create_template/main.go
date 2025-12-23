package main

import (
	"fmt"
	"log"

	"github.com/palemoky/lucky-day/internal/datasource"
)

func main() {
	// Create Excel template
	if err := datasource.CreateExcelTemplate("lottery_template.xlsx"); err != nil {
		log.Fatalf("Failed to create Excel template: %v", err)
	}
	fmt.Println("Excel template created successfully!")
}
