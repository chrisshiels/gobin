// 'true_test.go'.
// Chris Shiels.


package main


import (
    "bytes"
    "io/ioutil"
    "testing"
)


func Test_main(t *testing.T) {
    var tests = []struct {
        args []string
        stdinfilename string
        expectexitstatus int
        expectstdoutfilename string
        expectstderrfilename string
    }{
        {
            []string { "command" },
            "",
            0,
            "",
            "",
        },
    }

    for  _, test := range tests {
        var bytesstdin, bytesexpectstdout, bytesexpectstderr []byte
        var stdin, expectstdout, expectstderr *bytes.Buffer
        var stdout, stderr *bytes.Buffer
        var err error

        if test.stdinfilename != "" {
            bytesstdin, err = ioutil.ReadFile(test.stdinfilename)
            if err != nil {
                t.Errorf("true_test: %s", err)
                continue
            }
        }
        stdin = bytes.NewBuffer(bytesstdin)

        if test.expectstdoutfilename != "" {
            bytesexpectstdout, err = ioutil.ReadFile(test.expectstdoutfilename)
            if err != nil {
                t.Errorf("true_test: %s", err)
                continue
            }
        }
        expectstdout = bytes.NewBuffer(bytesexpectstdout)

        if test.expectstderrfilename != "" {
            bytesexpectstderr, err = ioutil.ReadFile(test.expectstderrfilename)
            if err != nil {
                t.Errorf("true_test: %s", err)
                continue
            }
        }
        expectstderr = bytes.NewBuffer(bytesexpectstderr)

        stdout = new(bytes.Buffer)
        stderr = new(bytes.Buffer)

        exitstatus := _main(stdin, stdout, stderr, test.args)

        if test.expectexitstatus != exitstatus ||
           expectstdout.String() != stdout.String() ||
           expectstderr.String() != stderr.String() {
            t.Errorf("%v = (%d, %q, %q), want (%d, %s: %q, %s: %q)",
                     test.args,
                     exitstatus,
                     stdout.String(),
                     stderr.String(),
                     test.expectexitstatus,
                     test.expectstdoutfilename,
                     expectstdout.String(),
                     test.expectstderrfilename,
                     expectstderr.String())

        }
    }
}
