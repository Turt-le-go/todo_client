package main

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"

	"todo_client/src/utils"
	"todo_client/src/todo"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const(
	connPort int = 8080
	connHost string = "localhost"

	timeLayout string = "02.01.06"
)

var reader = bufio.NewReader(os.Stdin)

func main(){
	fmt.Println("Client started")

	conn, err := grpc.Dial(fmt.Sprintf("%v:%d", connHost, connPort), grpc.WithInsecure())
	utils.Check(err)
	defer conn.Close()

	c := todo.NewToDoServiceClient(conn)
	
	commandLoop(c)
}

func commandLoop(c todo.ToDoServiceClient){
	for {
		input := getInput("-> ")
		switch input {
		case "add":
			addTask(c)
		case "list":
			printAllTasks(c)
		case "q":
			fmt.Println("Bye.")
			return
		default:
			fmt.Println("Wrong command.")
		}
	}
}

func getInput(prompt string) string {
		fmt.Print(prompt)
		input,_ := reader.ReadString('\n')
		return strings.Trim(input," \n")
}

func addTask(c todo.ToDoServiceClient){
	title := getInput("Task title: ")
	description := getInput("Task description: ")
	timeString := getInput("Deadline: ")
	deadline, err := parseTime(timeString)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Can't convert %v to time.\n", timeString)
		fmt.Printf("Use this %v format.\n", timeLayout)
		timeString = getInput("Deadline: ")
		deadline, err = parseTime(timeString)
		if err != nil {
			fmt.Printf("Can't convert %v to time.\n", timeString)
			fmt.Println("Failed to create task")
			return
		}
	}
	task := todo.TaskMessage{
		Title: title,
		Description: description,
		Deadline: deadline,
		CreatedAt: time.Now().Unix(),
	}
	res, err := c.AddTask(context.Background(), &task)
	utils.Check(err)
	fmt.Println(res.Text)
}

func parseTime(timeString string) (int64, error) {
	t, err := time.Parse(timeLayout, timeString)
	return t.Unix(), err
}

func printAllTasks(c todo.ToDoServiceClient) {
	tasksList, err := c.ListTasks(context.Background(), &todo.Empty{})
	utils.Check(err)
	terminalWigth,_ := utils.GetSize()
	delimiter := strings.Repeat("#", terminalWigth)
	fmt.Println(delimiter)
	for _,task := range tasksList.List {
		fmt.Printf("Title: %v\n", task.Title)
		fmt.Printf("Description: %v\n", task.Description)
		fmt.Printf("CreatedAt: %v\n", time.Unix(int64(task.CreatedAt), 0).Format(timeLayout))
		fmt.Printf("Deadline: %v\n", time.Unix(int64(task.Deadline), 0).Format(timeLayout))
		fmt.Println(delimiter)
	}
}
