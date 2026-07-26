package output

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"time"
)

var (
	Success = color.New(color.FgGreen, color.Bold).SprintFunc()
	Failure = color.New(color.FgRed, color.Bold).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Dim     = color.New(color.FgHiBlack).SprintFunc()
)

func Check(label string) {
	fmt.Printf("%s %s\n", Success("✓"), label)
}

func Cross(label string) {
	fmt.Printf("%s %s\n", Failure("✗"), label)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Info("ℹ"), fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Warning("⚠"), fmt.Sprintf(format, args...))
}

func Dimf(format string, args ...interface{}) {
	fmt.Println(Dim(fmt.Sprintf(format, args...)))
}

func NewSpinner(text string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + text
	s.Writer = nil
	return s
}
