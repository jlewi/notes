build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

# use port 8000 because we will use 8080 for the agenttraces server
run: build
	./gowasm --port=8000