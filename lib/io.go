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

	fmt.Fprintf(ios.errWriter, "%s: %v\n", prefix, err)
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

func (pr *PromptReader) ReadChoose(hdr string, choices []string, prompt string, defValue int, canDefault bool) (int, bool) {
	minChoice := 1
	maxChoice := len(choices)

	fmtPrompt := pr.formatChoosePrompt(prompt, minChoice, maxChoice, defValue, canDefault)

	for {
		if len(hdr) > 0 {
			fmt.Fprintf(pr.outWriter, "\n%s:\n\n", hdr)
		}

		for cidx, choice := range choices {
			fmt.Fprintf(pr.outWriter, "  %d - %s\n", cidx+1, choice)
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
				return choice - 1, false
			}

			continue
		}

		if canDefault {
			return defValue, true
		}
	}

	return 0, false
}

func (pr *PromptReader) ReadYN(prompt string, defValue, canDefault bool) (bool, bool) {
	fmtPrompt := pr.formatYNPrompt(prompt, defValue, canDefault)

	for {
		value, err := pr.promptReadTrimmedString(fmtPrompt)

		if err != nil {
			continue
		}

		if len(value) > 0 {
			switch strings.ToLower(value) {
			case "n", "no":
				return false, false

			case "y", "yes":
				return true, false

			default:
				continue
			}
		}

		if canDefault {
			return defValue, true
		}
	}

	return false, false
}

//-----------------------------------------------------------------------------

func (pr *PromptReader) formatChoosePrompt(prompt string, minChoice, maxChoice, defValue int, canDefault bool) string {
	if !canDefault {
		return fmt.Sprintf("\n%s (%d..%d): ", prompt, minChoice, maxChoice)
	}

	return fmt.Sprintf("\n%s (%d..%d): [%d] ", prompt, minChoice, maxChoice, defValue+1)
}

func (pr *PromptReader) formatYNPrompt(prompt string, defValue, canDefault bool) string {
	if !canDefault {
		return fmt.Sprintf("\n%s (Y/N)? ", prompt)
	}

	if defValue {
		return fmt.Sprintf("\n%s (Y/N)? [Y] ", prompt)
	}

	return fmt.Sprintf("\n%s (Y/N)? [N] ", prompt)
}

func (pr *PromptReader) promptReadTrimmedString(prompt string) (string, error) {
	fmt.Fprint(pr.outWriter, prompt)

	input, err := pr.inReader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(input, "\n"), nil
}
