package itlogs

import (
	"bytes"
	"log"
	"runtime"
	"strings"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	test := func(flags []LogMsgFlag) func(t *testing.T) {
		return func(t *testing.T) {
			var (
				b1 = &bytes.Buffer{}
				b2 = &bytes.Buffer{}
			)

			flag := 0
			for _, f := range flags {
				flag = flag | int(f)
			}

			stdLogger := log.New(b1, Trace.String()+" ", flag)
			itlogsLogger := NewLogger(b2, flags, Trace)

			stdLogger.Print("test")
			itlogsLogger.Log(Trace, "test")

			if b1.String() != b2.String() {
				t.Fatalf("b1 (%s) not equals b2 (%s), flags: %+v", b1.String(), b2.String(), flags)
			}
		}
	}

	flags := []LogMsgFlag{Ldate, Ltime, LUTC, Lmsgprefix, LstdFlags}
	for _, f := range flags {
		t.Run(f.String(), test([]LogMsgFlag{f}))
	}

	t.Run("Flags combination", test([]LogMsgFlag{Ldate, Ltime, Lmicroseconds, Lmsgprefix}))

	t.Run("Default logger setter", func(t *testing.T) {
		b := &bytes.Buffer{}
		l1 := NewDefaultLogger(b, Debug)
		SetDefaultLogger(l1)
		l2 := GetDefaultLogger()
		l2.Log(Debug, "test")

		if !strings.Contains(b.String(), "test") {
			t.Fatalf("b (%s) not contains 'test', level: %s", b.String(), Trace.String())
		}
	})

	t.Run("With file", func(t *testing.T) {
		var (
			b1 = &bytes.Buffer{}
			b2 = &bytes.Buffer{}
			l1 = NewLogger(b1, []LogMsgFlag{Llongfile}, Debug)
			l2 = NewLogger(b2, []LogMsgFlag{Lshortfile}, Debug)
		)

		l1.Log(Debug, "")
		l2.Log(Debug, "")

		_, fName, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatalf("caller filename is unknown")
		}

		shortFName := fName
		for i := len(fName) - 1; i > 0; i-- {
			if fName[i] == '/' {
				shortFName = fName[i+1:]
				break
			}
		}

		if !strings.Contains(b1.String(), fName) {
			t.Fatalf("b1 (%s) not contains full file name (%s)", b1.String(), fName)
		}

		if !strings.Contains(b2.String(), " "+shortFName) {
			t.Fatalf("b2 (%s) not contains short file name (%s)", b2.String(), shortFName)
		}
	})

	t.Run("Ignore small log level", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := NewDefaultLogger(b, Debug)
		l.Trace("test")

		if b.String() != "" {
			t.Fatalf("b (%s) is not empty", b.String())
		}
	})

	t.Run("Test levels functions", func(*testing.T) {
		test := func(b *bytes.Buffer, lvl LogLevel) {
			if !strings.Contains(b.String(), "test") {
				t.Fatalf("b (%s) not contains 'test', level: %s", b.String(), lvl.String())
			}
			b.Reset()
		}

		b := &bytes.Buffer{}
		l := NewDefaultLogger(b, Trace)
		s := "test"
		l.Trace(s)
		test(b, Trace)
		l.Debug(s)
		test(b, Debug)
		l.Info(s)
		test(b, Info)
		l.Warn(s)
		test(b, Warning)
		l.Error(s)
		test(b, Error)
		l.Fatal(s)
		test(b, Fatal)

	})
}
