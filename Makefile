run_go_dev:
	cd app/cmd/uptime/; go run main.go

run_web_dev:
	cd frontend; npm run start

build_go:
	cd app/cmd/uptime/; go build

build_web:
	cd frontend; npm run build; cp -r ./build ../app/cmd/uptime

build:
	cd frontend; npm run build; cp -r ./build ../app/cmd/uptime
	go test ./... -cover
	cd app/cmd/uptime/; go build

build_linux:
	cd frontend; npm run build
	go test ./... -cover
	cd app/cmd/uptime/; GOOS=linux GOARCH=amd64 rice embed-go; GOOS=linux GOARCH=amd64 go build

test:
	go test ./... -cover

load_test:
	k6 run --summary-trend-stats="med,p(95),p(99.9)" load_testing/script.js

benchmark:
	cd load_testing; plow ${UPTIME_HOST}/API/v1/services -c 100 -n 10000 -T 'application/json' -m GET -H "Authorization: Bearer ${UPTIME_TOKEN}"

run_test_server:
	cd app/cmd/testserver/; go run main.go --port 1315
