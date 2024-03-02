build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

run: build
	./gowasm --port=8000