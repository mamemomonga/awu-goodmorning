bin/goodmorning:
	go build -o $@ ./cmd/

clean:
	rm -rf bin