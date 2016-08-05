// 'tail.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
    "time"
)


const exitsuccess = 0
const exitfailure = 1


type tailer struct {
    stdout io.Writer
    filename string
    reader *bufio.Reader
}


func newtailer(stdout io.Writer, filename string, reader io.Reader) *tailer {
    return &tailer{ stdout: stdout,
                    filename: filename,
                    reader: bufio.NewReader(reader) }
}


func (t *tailer) tail(nlines int) error {
    linen := 0
    lines := make([]string, nlines)

    var s string
    var err error
    for s, err = t.reader.ReadString('\n');
        err == nil;
        s, err = t.reader.ReadString('\n') {
        lines[linen % nlines] = s[0:len(s) - 1]
        linen++
    }

    if err != io.EOF {
        return err
    }

    if linen >= nlines {
        for n := linen % nlines; n < nlines; n++ {
            fmt.Fprintln(t.stdout, lines[n])
        }
    }
    for n := 0; n < linen % nlines; n++ {
        fmt.Fprintln(t.stdout, lines[n])
    }

    return nil
}


func (t *tailer) canfollow() (ret bool, err error) {
    if _, err := t.reader.ReadByte(); err != nil {
        if err == io.EOF {
            err = nil
        }
        return false, err
    }

    if err := t.reader.UnreadByte(); err != nil {
        return false, err
    }

    return true, nil
}


func (t *tailer) follow() error {
    var s string
    var err error
    for s, err = t.reader.ReadString('\n');
        err == nil;
        s, err = t.reader.ReadString('\n') {
        fmt.Fprintln(t.stdout, s[0:len(s) - 1])
    }

    if err == io.EOF {
        err = nil
    }
    return err
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {

    flagset := flag.NewFlagSet(args[0], flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stdout, "Usage:  tail [ -n n ] [ -f ] [ file ... ]")
        flag.PrintDefaults()
    }
    flagf := flagset.Bool("f",
                          false,
                          "Follow appended output by file descriptor")
    flagn := flagset.Int("n",
                         10,
                         "Output the last n lines")

    // Note flagset.Parse() will also handle '-h' and '--help' and will exit
    // with exit status 2.
    flagset.Parse(args[1:])


    var tailers []*tailer

    if len(flagset.Args()) == 0 {
        tailer := newtailer(stdout, "-", stdin)
        tailer.tail(*flagn)
        tailers = append(tailers, tailer)
    } else {
        for i, filename := range flagset.Args() {
            file, err := os.Open(filename)
            if err != nil {
                fmt.Fprintf(stderr, "tail: %s\n", err)
                continue
            }

            if len(flagset.Args()) > 1 {
                if i != 0 {
                    fmt.Fprintln(stdout)
                }
                fmt.Fprintf(stdout, "==> %s <==\n", filename)
            }

            tailer := newtailer(stdout, filename, file)
            err = tailer.tail(*flagn)
            if err != nil {
                fmt.Fprintf(stderr, "tail: %s\n", err)
                continue
            }
            tailers = append(tailers, tailer)
        }
    }

    if len(tailers) == 0 || !*flagf {
        return 0
    }

    // Update tailers[0] so the most recent tailer as at index 0.
    tailers[0], tailers[len(tailers) - 1] =
        tailers[len(tailers) - 1], tailers[0]

    for {
        for i, tailer := range tailers {

            canfollow, err := tailer.canfollow()
            if err != nil {
                fmt.Fprintf(stderr, "tail: %s\n", err)
                continue
            }

            if !canfollow {
                continue
            }

            if len(flagset.Args()) > 1 {
                if i != 0 {
                    fmt.Fprintln(stdout)
                    fmt.Fprintf(stdout, "==> %s <==\n", tailer.filename)
                }
            }

            err = tailer.follow()
            if err != nil {
                fmt.Fprintf(stderr, "tail: %s\n", err)
                continue
            }

            if i != 0 {
                // Update tailers[0] so the most recent tailer as at index 0.
                tailers[0], tailers[i] = tailers[i], tailers[0]
            }
        }
        time.Sleep(time.Second)
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
