run_go_dev:
	cd app/cmd/uptime/; swag init --pd true; go run main.go

run_web_dev:
	cd frontend; npm run start

build_go:
	cd app/cmd/uptime/; swag init --pd true; rice embed-go; go build

build_web:
	cd frontend; npm run build

build:
	cd frontend; npm run build
	go test ./... -cover
	cd app/cmd/uptime/; swag init --pd true; rice embed-go; go build