release:
	go get
	go get -u -v github.com/laher/goxc
	goxc -tasks='xc archive' -bc 'linux windows darwin' -d .
clean:
	rm -rf snapshot
