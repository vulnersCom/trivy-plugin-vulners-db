APP?=trivy-vulners-db

build:
	go build -buildvcs=false -o ./bin/${APP} ./cmd