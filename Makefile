all:
	go build  -gcflags "-N -l"
	go install

testbuild:
	go test -c -gcflags "-N -l" -v

clean:
	rm -f *~  webiperf *.o


