package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type powerShell struct {
	shell  *exec.Cmd
	stdin  io.Writer
	stdout bufferedStreamReader
	stderr bufferedStreamReader
}

func (s *powerShell) open() error {
	s.shell = exec.Command(`powerShell.exe`, `-NoExit`, `-Command`, `-`)
	s.shell.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdin, err := s.shell.StdinPipe()
	if err != nil {
		return err
	}
	s.stdin = stdin

	stdout, err := s.shell.StdoutPipe()
	if err != nil {
		return err
	}
	s.stdout = bufferedStreamReader{
		reader:          stdout,
		streamSeparator: powerShellStreamSeparator + powerShellStreamNewline,
	}

	stderr, err := s.shell.StderrPipe()
	if err != nil {
		return err
	}
	s.stderr = bufferedStreamReader{
		reader:          stderr,
		streamSeparator: powerShellStreamSeparator + powerShellStreamNewline,
	}

	err = s.shell.Start()
	if err != nil {
		return err
	}

	s.syncExecute(powerShellInitCommand)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.close()
	}()

	return nil
}

func (s *powerShell) close() {
	s.executeLine(`exit`)
	s.shell = nil
	s.stdin = nil
	s.stdout.reader = nil
	s.stderr.reader = nil
}

func (s *powerShell) syncExecuteCount(command string) (int64, error) {
	err := s.syncExecute(`@(` + command + `).Count`)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(s.stdout.buffer), 10, 64)
}

func (s *powerShell) syncExecute(command string) error {
	err := s.executeLine(
		fmt.Sprintf(
			"%s ; echo '%s' ; [Console]::Error.WriteLine('%s')",
			command, powerShellStreamSeparator, powerShellStreamSeparator))
	if err != nil {
		return fmt.Errorf("Unable to write to powershell stream: %v", err)
	}

	if err := s.wait(); err != nil {
		return fmt.Errorf("Unable to read from powershell stream: %v", err)
	}

	return nil
}

func (s *powerShell) executeLine(command string) error {
	return s.execute(command + powerShellStreamNewline)
}

func (s *powerShell) execute(command string) error {
	_, err := s.stdin.Write([]byte(command))
	return err
}

func (s *powerShell) wait() error {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2)
	var err error
	completedFunc := func(readError error) {
		if readError != nil {
			err = readError
		}
		waitGroup.Done()
	}
	s.stdout.asyncRead(completedFunc)
	s.stderr.asyncRead(completedFunc)

	waitGroup.Wait()

	return err
}

type bufferedStreamReader struct {
	reader          io.Reader
	streamSeparator string
	buffer          string
}

func (b *bufferedStreamReader) asyncRead(completedFunc func(err error)) {
	go func() {
		completedFunc(b.read())
	}()
}

func (b *bufferedStreamReader) read() error {
	readBuffer := make([]byte, 32)
	output := ""

	for {
		read, err := b.reader.Read(readBuffer)
		if err != nil {
			return err
		}

		output += string(readBuffer[:read])

		if strings.HasSuffix(output, b.streamSeparator) {
			break
		}
	}

	b.buffer = strings.TrimSuffix(output, b.streamSeparator)

	return nil
}
