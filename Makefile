all:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build
	zip -9 /tmp/go2.zip dict.txt freebee
