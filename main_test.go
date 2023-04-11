package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()

	// Exit with the code that the tests returned
	os.Exit(code)
}

func Test_readUserInput(t *testing.T) {
	type args struct {
		in       io.Reader
		doneChan chan bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "q",
			args: args{
				in:       strings.NewReader("q"),
				doneChan: make(chan bool),
			},
		},
		{
			name: "Q",
			args: args{
				in:       strings.NewReader("Q"),
				doneChan: make(chan bool),
			},
		}, {
			name: "number",
			args: args{
				in:       strings.NewReader("7\nq"),
				doneChan: make(chan bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go readUserInput(tt.args.in, tt.args.doneChan)
			<-tt.args.doneChan
		})
	}
}

func Test_checkNumbers(t *testing.T) {
	type args struct {
		scanner *bufio.Scanner
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantBool bool
	}{
		{
			name: "q",
			args: args{
				scanner: bufio.NewScanner(strings.NewReader("q")),
			},
			want:     "",
			wantBool: true,
		},
		{
			name: "Q",
			args: args{
				scanner: bufio.NewScanner(strings.NewReader("Q")),
			},
			want:     "",
			wantBool: true,
		},
		{
			name: "number",
			args: args{
				scanner: bufio.NewScanner(strings.NewReader("f")),
			},
			want:     "Please enter a whole number!",
			wantBool: false,
		},
		{
			name: "valid number",
			args: args{
				scanner: bufio.NewScanner(strings.NewReader("7")),
			},
			want:     "7 is a prime number!",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, book := checkNumbers(tt.args.scanner); got != tt.want || book != tt.wantBool {
				t.Errorf("checkNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntro(t *testing.T) {
	// Redirect stdout to a buffer so we can check the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	intro()

	// Reset stdout
	w.Close()
	os.Stdout = old

	// Check the output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	expected := "Is it Prime?\n------------\nEnter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit.\n-> "
	if buf.String() != expected {
		t.Log(buf.String())
		t.Errorf("Intro did not return correct string. Expected %q but got %q", expected, buf.String())
	}
}

func TestPrompt(t *testing.T) {
	// Redirect stdout to a buffer so we can check the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	prompt()

	// Reset stdout
	w.Close()
	os.Stdout = old

	// Check the output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	expected := "-> "
	if buf.String() != expected {
		t.Errorf("Prompt did not return correct string. Expected %q but got %q", expected, buf.String())
	}
}

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}
