swagger:
	swag init  --parseDependency --parseInternal

build-linux:
	CC=x86_64-unknown-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o meego_meeting_plugin . && upx -5 meego_meeting_plugin