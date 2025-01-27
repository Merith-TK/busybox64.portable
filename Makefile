default: setup
	go build ./ -o ./busybox64.portable
	export PATH=""
	./busybox64.portable

windows: setup
	GOOS=windows go build -o ./busybox64.portable.exe ./
	export PATH=""
	./busybox64.portable.exe

setup:
	go generate ./...