// 'head.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


type header struct {
    stdout io.Writer
    filename string
    reader *bufio.Reader
}


func newheader(stdout io.Writer, filename string, reader io.Reader) *header {
    return &header{ stdout: stdout,
                    filename: filename,
                    reader: bufio.NewReader(reader) }
}


func (h *header) head(nlines int) error {
    linen := 0

    var s string
    var err error
    for s, err = h.reader.ReadString('\n');
        err == nil;
        s, err = h.reader.ReadString('\n') {
        fmt.Fprintln(h.stdout, s[0:len(s) - 1])
        linen++
        if linen == nlines {
            break
        }
    }

    if err != io.EOF {
        return err
    }
    return nil
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {

    flagset := flag.NewFlagSet(args[0], flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stdout, "Usage:  head [ -n n ] [ file ... ]")
        flag.PrintDefaults()
    }
    flagn := flagset.Int("n",
                         10,
                         "Output the last n lines")

    // Note flagset.Parse() will also handle '-h' and '--help' and will exit
    // with exit status 2.
    flagset.Parse(args[1:])


    if len(flagset.Args()) == 0 {
        header := newheader(stdout, "-", stdin)
        header.head(*flagn)
    } else {
        for i, filename := range flagset.Args() {
            file, err := os.Open(filename)
            if err != nil {
                fmt.Fprintf(stderr, "head: %s\n", err)
                continue
            }

            if len(flagset.Args()) > 1 {
                if i != 0 {
                    fmt.Fprintln(stdout)
                }
                fmt.Fprintf(stdout, "==> %s <==\n", filename)
            }

            header := newheader(stdout, filename, file)
            err = header.head(*flagn)
            if err != nil {
                fmt.Fprintf(stderr, "head: %s\n", err)
                continue
            }
        }
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
