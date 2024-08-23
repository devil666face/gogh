package view

import (
	"fmt"
	"gogh/internal/gogh"
	"gogh/internal/view/options"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

const (
	Emoji    = "🪣"
	NotFound = "⚠️ command not found"
)

var ErrorFunc func(err error) = func(err error) {
	fmt.Printf("⚠️ " + err.Error() + "\r\n")
}

var boolFormat = func(c bool) string {
	if c {
		return "+"
	}
	return " "
}

var NotFoundFunc func() = func() {
	fmt.Println(NotFound)
}

var (
	Title = fmt.Sprintf("%s >> ", Emoji)
)

const (
	Store      = "store"
	Settings   = "settings"
	Show       = "show"
	Delete     = "delete"
	Upload     = "upload"
	NoCompress = "nocompress"
	Compress   = "compress"
	Download   = "download"
	Token      = "token"
	Set        = "set"
	Exit       = "exit"
)

type View struct {
	prompt *prompt.Prompt
	gogh   *gogh.Gogh
}

func (v *View) executor(in string) {
	args := strings.Fields(in)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case Store:
		NewPrompt(
			v.storeExecutor,
			v.storeCompleter,
			fmt.Sprintf("%s%s >> ", Title, "🗄️ "+Store),
		).Run()
	case Settings:
		NewPrompt(
			v.settingsExecutor,
			v.settingsCompleter,
			fmt.Sprintf("%s%s >> ", Title, "⚙️ "+Settings),
		).Run()
	case Exit:
		os.Exit(0)
	default:
		NotFoundFunc()
	}
}

func (v *View) completer(d prompt.Document) []prompt.Suggest {
	complete := []prompt.Suggest{
		{Text: Store, Description: "file management"},
		{Text: Settings, Description: "configure"},
		{Text: Exit, Description: "close panel"},
	}
	// Remove second postition complete
	for _, c := range complete {
		if HasPrefix(d, c.Text) {
			complete = []prompt.Suggest{}
		}
	}
	return prompt.FilterContains(complete, d.GetWordBeforeCursor(), true)
}

func New(_gogh *gogh.Gogh) *View {
	v := &View{
		gogh: _gogh,
	}
	v.prompt = NewPrompt(
		v.executor,
		v.completer,
		Title,
	)
	return v
}

func NewPrompt(executor func(string), completer func(prompt.Document) []prompt.Suggest, title string) *prompt.Prompt {
	return prompt.New(
		executor,
		completer,
		options.PromptOptions(
			title,
			Exit,
		)...,
	)
}

func (v *View) Run() {
	v.prompt.Run()
}

func HasPrefix(d prompt.Document, prefix string) bool {
	return strings.HasPrefix(strings.TrimSpace(d.TextBeforeCursor()), prefix)
}
