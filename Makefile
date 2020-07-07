gelato: *.go
	gofmt -w *.go
	go build -o gelato *.go
