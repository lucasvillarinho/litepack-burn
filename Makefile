
.PHONY: run-app
run-app:
	@echo "Running main.go..."
	go run ./cmd/api/main.go

.PHONY: run-attack
run-attack:
	@echo "Running atack.go..."
	go run ./cmd/attack/main.go