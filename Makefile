default:
	go build ./ -o ./busybox64.portable
	export PATH=""
	./busybox64.portable

windows:
	GOOS=windows go build -o ./busybox64.portable.exe ./
	export PATH=""
	./busybox64.portable.exe