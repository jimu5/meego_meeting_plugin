swagger:
	swag init  --parseDependency --parseInternal

build-linux:
	GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o meego_meeting_plugin .

build-linux-with-upx:
	GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o meego_meeting_plugin . && upx -5 meego_meeting_plugin