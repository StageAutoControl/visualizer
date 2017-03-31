proto:
	protoc -I "dmx" --go_out="dmx" dmx/dmx.proto

build:
	go build -o visualizer .
