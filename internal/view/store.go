package view

import (
	"fmt"
	"gogh/pkg/fs"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

func (v *View) uploadAction() {
	var (
		file string
		opts []huh.Option[string]
	)
	files, err := fs.FilesInCurrentDir()
	if err != nil {
		ErrorFunc(err)
		return
	}
	for _, f := range files {
		opts = append(opts, huh.NewOption[string](f, f))
	}

	if err := huh.NewSelect[string]().
		Title("🗃 select a file").
		Options(opts...).
		Value(&file).Run(); err != nil {
		ErrorFunc(err)
		return
	}
	var upload func() = func() {
		if err := v.gogh.Upload(file); err != nil {
			ErrorFunc(err)
			return
		}
		fmt.Printf("👌 uploaded %s\n", file)
	}
	_ = spinner.New().Title(fmt.Sprintf("⏳ upload %s", file)).Action(upload).Run()
}

func (v *View) storeExecutor(in string) {
	args := strings.Fields(in)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case Upload:
		v.uploadAction()
	case Show:
		var (
			sb    strings.Builder
			w     = tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
			i     = 1
			title = func(s string) string {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).Render(s)
			}
		)
		fmt.Fprintf(w, "#\t%s\t%s", title("Filename"), title("Created"))
		for _, v := range v.gogh.Data.Filestore.Files {
			fmt.Fprintf(w, "\n%d\t%s\t%s", i, v.Filename, v.CreatedDate)
			i++
		}
		w.Flush()
		fmt.Println(
			lipgloss.NewStyle().
				// BorderStyle(lipgloss.RoundedBorder()).
				// BorderForeground(lipgloss.Color("63")).
				Padding(0, 1).
				Render(sb.String()),
			// Width(40).
		)
	case Exit:
		os.Exit(0)
	default:
		NotFoundFunc()
	}
}

func (v *View) storeCompleter(d prompt.Document) []prompt.Suggest {
	var complete = []prompt.Suggest{
		{Text: Upload, Description: "upload local file to store"},
		{Text: Download, Description: "download remote file"},
		{Text: Show, Description: "show all uploaded files"},
		{Text: Exit, Description: ""},
	}
	// remove second postition complete
	for _, c := range complete {
		if HasPrefix(d, c.Text) {
			complete = []prompt.Suggest{}
		}
	}
	return prompt.FilterContains(complete, d.GetWordBeforeCursor(), true)
}
