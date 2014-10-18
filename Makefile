all: willus

willus: src/*.go
	go build -o $@ src/*.go

clean:
	rm willus
