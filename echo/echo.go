// 'echo.go'.
// Chris Shiels.


package main


import (
    "bytes"
    "flag"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)


const exitsuccess = 0
const exitfailure = 1


func expandescapesequences(s string) (expandeds string, err error) {
    var reader *strings.Reader = strings.NewReader(s)
    var ch1, ch2, ch3, ch4 rune
    var buffer bytes.Buffer
    var r rune

    for {
        if ch1, _, err = reader.ReadRune(); err != nil {
            break
        }

        if ch1 != '\\' {
            r = ch1
        } else {
            if ch2, _, err = reader.ReadRune(); err != nil {
               break
            }

            switch ch2 {
                case '\\':
                    r = '\\'
                case 'a':
                    r = '\a'
                case 'b':
                    r = '\b'
                case 'c':
                    break
                case 'e':
                    r = '\x1b'
                case 'f':
                    r = '\f'
                case 'n':
                    r = '\n'
                case 'r':
                    r = '\r'
                case 't':
                    r = '\t'
                case 'v':
                    r = '\v'
                case '0', '1', '2', '3', '4', '5', '6', '7':
                    if ch3, _, err = reader.ReadRune(); err != nil {
                        break
                    }
                    if ch4, _, err = reader.ReadRune(); err != nil {
                        break
                    }

                    var v uint64
                    v, err =
                        strconv.ParseUint(fmt.Sprintf("%c%c%c", ch2, ch3, ch4),
                                          8, 8)
                    r = rune(v)

                case 'x':
                    if ch3, _, err = reader.ReadRune(); err != nil {
                        break
                    }
                    if ch4, _, err = reader.ReadRune(); err != nil {
                        break
                    }

                    var v uint64
                    v, err =
                        strconv.ParseUint(fmt.Sprintf("%c%c", ch3, ch4),
                                          16, 8)
                    r = rune(v)

                default:
                    if _, err = buffer.WriteRune('\\'); err != nil {
                        break
                    }
                    r = ch2
            }
        }

        if _, err = buffer.WriteRune(r); err != nil {
            break
        }
    }

    if (err == io.EOF) {
        err = nil
    }

    return buffer.String(), err
}


func echo(stdout io.Writer, flage bool, flagn bool, args []string) error {
    var err error

    s := strings.Join(args, " ")

    if (flage) {
        if s, err = expandescapesequences(s); err != nil {
            return err
        }
    }

    fmt.Fprint(stdout, s)
    if (!flagn) {
        fmt.Fprintln(stdout)
    }

    return nil
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {

    flagset := flag.NewFlagSet(args[0], flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stderr, "Usage:  echo [ -e ] [ -n ] [ args ... ]")
        flagset.PrintDefaults()
    }
    flage := flagset.Bool("e",
                          false,
                          "Enable interpretation of backslash escapes")
    flagn := flagset.Bool("n",
                          false,
                          "Do not output the trailing newline")

    // Note flagset.Parse() will also handle '-h' and '--help' and will exit
    // with exit status 2.
    flagset.Parse(args[1:])

    if err := echo(stdout, *flage, *flagn, flagset.Args()); err != nil {
        fmt.Fprintf(stderr, "echo: %s\n", err)
        return exitfailure
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
