# gobin

Because rewriting simple versions of standard UNIX utilities is a great way
to learn a new language...

    host$ export GOPATH=~/workspace
    host$ mkdir -p $GOPATH
    host$ cd $GOPATH

    host$ go get github.com/chrisshiels/gobin/...

    host$ ./bin/cat ./src/github.com/chrisshiels/gobin/cat/cat.go

    host$ ./bin/echo -n -e '\e[31mHello'

    host$ ./bin/find . -type f -name \*.go -print

    host$ ./bin/grep ^ro /etc/passwd

    host$ ./bin/head -n 5 /etc/group /etc/passwd

    host$ ./bin/ls bin src

    host$ ./bin/tac ./src/github.com/chrisshiels/gobin/cat/cat.go | ./bin/tac

    host$ ./bin/tail -f /var/log/cron

    host$ ./bin/true ; echo $?

    host$ ./bin/wc ./src/github.com/chrisshiels/gobin/*/*.go
