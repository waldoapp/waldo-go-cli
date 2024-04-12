package lib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type IOStreams struct {
	inReader  io.Reader
	outWriter io.Writer
	errWriter io.Writer
}

//-----------------------------------------------------------------------------

func NewIOStreams(in io.Reader, out, err io.Writer) *IOStreams {
	return &IOStreams{
		inReader:  in,
		outWriter: out,
		errWriter: err}
}

//-----------------------------------------------------------------------------

func (ios *IOStreams) EmitError(prefix string, err error) {
	fmt.Fprintf(ios.outWriter, "\n") // flush output

	fmt.Fprintf(ios.errWriter, "%v: %v\n", prefix, err)
}

func (ios *IOStreams) Print(a ...any) (n int, err error) {
	return fmt.Fprint(ios.outWriter, a...)
}

func (ios *IOStreams) Println(a ...any) (n int, err error) {
	return ios.Print(fmt.Sprintln(a...))
}

func (ios *IOStreams) Printf(format string, a ...any) (n int, err error) {
	return ios.Print(fmt.Sprintf(format, a...))
}

func (ios *IOStreams) PrintErr(a ...any) (n int, err error) {
	return fmt.Fprint(ios.errWriter, a...)
}

func (ios *IOStreams) PrintErrln(a ...any) (n int, err error) {
	return ios.PrintErr(fmt.Sprintln(a...))
}

func (ios *IOStreams) PrintErrf(format string, a ...any) (n int, err error) {
	return ios.PrintErr(fmt.Sprintf(format, a...))
}

//-----------------------------------------------------------------------------

type PromptReader struct {
	inReader  *bufio.Reader
	outWriter io.Writer
}

//-----------------------------------------------------------------------------

func (ios *IOStreams) PromptReader() *PromptReader {
	return &PromptReader{
		inReader:  bufio.NewReader(ios.inReader),
		outWriter: ios.outWriter}
}

//-----------------------------------------------------------------------------

func (pr *PromptReader) ReadChoose(hdr string, choices []string, prompt string) int {
	minChoice := 1
	maxChoice := len(choices)

	fmtPrompt := pr.formatChoosePrompt(prompt, minChoice, maxChoice)

	for {
		if len(hdr) > 0 {
			fmt.Fprintf(pr.outWriter, "\n%v:\n\n", hdr)
		}

		for cidx, choice := range choices {
			fmt.Fprintf(pr.outWriter, "  %v - %v\n", cidx+1, choice)
		}

		value, err := pr.promptReadTrimmedString(fmtPrompt)

		if err != nil {
			continue
		}

		var choice int

		if len(value) > 0 {
			choice, err = strconv.Atoi(value)

			if err != nil {
				continue
			}

			if choice >= minChoice && choice <= maxChoice {
				return choice - 1
			}

			continue
		}
	}

	return 0
}

func (pr *PromptReader) ReadYN(prompt string) bool {
	fmtPrompt := pr.formatYNPrompt(prompt)

	for {
		value, err := pr.promptReadTrimmedString(fmtPrompt)

		if err != nil {
			continue
		}

		if len(value) > 0 {
			switch strings.ToLower(value) {
			case "n", "no":
				return false

			case "y", "yes":
				return true

			default:
				continue
			}
		}
	}

	return false
}

//-----------------------------------------------------------------------------

func (pr *PromptReader) formatChoosePrompt(prompt string, minChoice, maxChoice int) string {
	return fmt.Sprintf("\n%v (%d..%d): ", prompt, minChoice, maxChoice)
}

func (pr *PromptReader) formatYNPrompt(prompt string) string {
	return fmt.Sprintf("\n%v (Y/N)? ", prompt)
}

func (pr *PromptReader) promptReadTrimmedString(prompt string) (string, error) {
	fmt.Fprint(pr.outWriter, prompt)

	input, err := pr.inReader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(input, "\n"), nil
}
