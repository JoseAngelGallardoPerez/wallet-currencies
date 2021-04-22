// Package makes possible to execute task written in go in context of app
package tasks

import (
	"fmt"
	"log"

	"go.uber.org/dig"
)

// Available tasks
var tasks = map[string]interface{}{
	"populateDb":  populateDb,
	"updateRates": updateRates,
}

// Executes passed task
func Execute(name *string, container *dig.Container) {
	if task, ok := tasks[*name]; ok {
		err := container.Invoke(task)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Task %v was not found\n", *name)
	}
}
