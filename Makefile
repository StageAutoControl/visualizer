build:
	go build -o visualizer .

start: build
	./visualizer server
