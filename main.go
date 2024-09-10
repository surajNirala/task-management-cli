package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

const taskFile = "task.json"

var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var reset = "\033[0m"

func main() {
	for {
		fmt.Println(string(yellow), "Task List", string(reset))
		fmt.Println(string(yellow), "1. Add Task", string(reset))
		fmt.Println(string(yellow), "2. View Task", string(reset))
		fmt.Println(string(yellow), "3. Delete Task", string(reset))
		fmt.Println(string(yellow), "4. Mark Task as Complete", string(reset))
		fmt.Println(string(yellow), "5. Exit", string(reset))
		fmt.Println(string(yellow), "Enter Your Choice: ", string(reset))

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			AddTask()
		case "2":
			ViewTasks()
		case "3":
			DeleteTask()
		case "4":
			CompleteTask()
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice Please try again.")
		}
	}

}

func LoadTasks() ([]Task, error) {
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		return []Task{}, nil
	}
	data, err := os.ReadFile(taskFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file : %w", err)
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks from JSON : %w", err)
	}
	return tasks, nil
}

func AddTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter task description.")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	// taskFile
	newtask := Task{
		Description: description,
		Completed:   false,
	}
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Printf("Error Loading tasks: %v\n", err)
		return
	}
	tasks = append(tasks, newtask)
	err = SaveTask(tasks)
	if err != nil {
		fmt.Printf("Error Saving Tasks: %v\n", err)
		return
	}
	fmt.Println(string(green), "Task Created successfully.", string(reset))
}

func SaveTask(tasks []Task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}
	err = os.WriteFile(taskFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks to file :%w", err)
	}
	return nil
}

func ViewTasks() {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Printf("Error Loading tasks %v", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	// fmt.Printf("%-5s %-30s %s\n", "ID", "Description", "Completed")
	fmt.Printf("%s%-5s %-30s %s%s\n", yellow, "ID", "Description", "Completed", reset)
	fmt.Println(string(yellow), "-------------------------------------------------", string(reset))
	for i, task := range tasks {
		fmt.Printf("%-5s %-30s %t\n", i+1, task.Description, task.Completed)
	}
	fmt.Println(string(yellow), "-------------------------------------------------", string(reset))
}

func CompleteTask() {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Printf("Error Loading tasks %v", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter task number for update.")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	tasks[taskNum-1].Completed = true
	err = SaveTask(tasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		return
	}
	fmt.Println(string(green), "Task Updated successfully.", string(reset))
}

func DeleteTask() {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Printf("Error Loading tasks %v", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	// fmt.Println("Enter task description.")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter task number for delete.")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	tasks = append(tasks[:taskNum-1], tasks[taskNum:]...)
	err = SaveTask(tasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		return
	}
	fmt.Println(string(red), "Task deleted successfully.", string(reset))
}
