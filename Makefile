default: run

clear:
	rm todo

run:
	go build -o ./todo src/main.go
	./todo
