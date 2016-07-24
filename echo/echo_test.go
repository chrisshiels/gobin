// 'echo_test.go'.
// Chris Shiels.


package main


import (
    "bytes"
    "testing"
)


func TestEcho(t *testing.T) {
    var tests = []struct {
        args []string
        expectexitstatus int
        expectstdout string
        expectstderr string
    }{
        {
            []string { "command" },
            0, "\n", "",
        },
        {
            []string { "command", "-n" },
            0, "", "",
        },
        {
            []string { "command", "a", "b", "c" },
            0, "a b c\n", "",
        },
        {
            []string { "command", "-n", "a", "b", "c" },
            0, "a b c", "",
        },
        {
            []string { "command", "-e", `a\nb\nc` },
            0, "a\nb\nc\n", "",
        },
        {
            []string { "command", "-n", "-e", `a\nb\nc` },
            0, "a\nb\nc", "",
        },
        {
            []string { "command", "-e", `\e[31mHello\e[0m` },
            0, "\x1b[31mHello\x1b[0m\n", "",
        },
        {
            []string { "command", "-n", "-e", `\e[31mHello\e[0m` },
            0, "\x1b[31mHello\x1b[0m", "",
        },
    }

    for _, test := range tests {
        var stdout bytes.Buffer
        var stderr bytes.Buffer

        exitstatus := maindo(nil, &stdout, &stderr, test.args)

        if test.expectexitstatus != exitstatus ||
           test.expectstdout != stdout.String() ||
           test.expectstderr != stderr.String() {
            t.Errorf("%v = (%d, %q, %q), want (%d, %q, %q)",
                     test.args,
                     exitstatus,
                     stdout.String(),
                     stderr.String(),
                     test.expectexitstatus,
                     test.expectstdout,
                     test.expectstderr)
        }
    }
}
