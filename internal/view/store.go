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

const dateFormat = "Jan 02 15:04"

var (
	base16 *huh.Theme = huh.ThemeBase16()
)

func RunSelect(
	title string,
	opts []huh.Option[string],
	value *string,
) error {
	s := huh.NewSelect[string]().
		Title("select file").
		Options(opts...).
		Value(value).
		WithTheme(base16)
	return huh.NewForm(huh.NewGroup(s)).Run()
}

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
	if err := RunSelect("upload file", opts, &file); err != nil {
		ErrorFunc(err)
		return
	}
	var upload func() = func() {
		if err := v.gogh.Upload(file); err != nil {
			ErrorFunc(err)
			return
		}
		fmt.Printf("👌 uploaded %s\r\n", file)
	}
	_ = spinner.New().Title(fmt.Sprintf("⏳ upload %s", file)).Action(upload).Run()
}

func (v *View) showAction() {
	var (
		sb strings.Builder
		w  = tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
	)
	fmt.Fprintf(w, "#\t%s\t%s\t%s", "Name", "Created", "Size")
	for i, v := range v.gogh.Data.Filestore.FilesSlice() {
		fmt.Fprintf(w, "\n%d\t%s\t%s\t%s", i+1, v.Filename, v.CreatedDate.Format(dateFormat), v.FormatSize())
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
}

func (v *View) deleteAction() {
	var (
		file string
		opts []huh.Option[string]
	)
	for k, f := range v.gogh.Data.Filestore.Files {
		opts = append(opts, huh.NewOption[string](f.Filename, k))
	}
	if err := RunSelect("delete file", opts, &file); err != nil {
		ErrorFunc(err)
		return
	}
	fmt.Println(file)
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
		v.showAction()
	case Delete:
		v.deleteAction()
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
		{Text: Delete, Description: "delete file from store"},
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
