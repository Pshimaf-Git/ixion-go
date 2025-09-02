OUTPUT = .out/ixion

MAIN_FILE = cmd/main.go

build:
	go build -o $(OUTPUT) $(MAIN_FILE)

run:
	go run $(MAIN_FILE)

build-run: build
	./$(OUTPUT)

clean:
	rm -f $(OUTPUT)

help:
	@echo "Commands:"
	@echo "  make build     - build app in $(OUTPUT)"
	@echo "  make run       - run app"
	@echo "  make build-run - build and run app"
	@echo "  make clean     - delete result of compilation"
	@echo "  make help      "

.DEFAULT_GOAL := build
