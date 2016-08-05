// 'true.go'.
// Chris Shiels.


package main


import (
    "io"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
