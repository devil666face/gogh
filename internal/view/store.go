package view

import (
	"fmt"
	"gogh/pkg/fs"
	"os"
	"strconv"
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
		Title(title).
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
	if err := RunSelect("upload", opts, &file); err != nil {
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
	if len(v.gogh.Data.Filestore.Files) == 0 {
		ErrorFunc(fmt.Errorf("no uploaded data"))
		return
	}
	fmt.Fprintf(w, "#\t%s\t%s\t%s", "Name", "Created", "Size")
	for i, v := range v.gogh.Data.Filestore.Files {
		fmt.Fprintf(w, "\n%d\t%s\t%s\t%s", i+1, v.Filename, v.CreatedDate.Format(dateFormat), v.FormatSize())
	}
	w.Flush()
	fmt.Println(
		lipgloss.NewStyle().
			Padding(0, 1).
			Render(sb.String()),
	)
}

func (v *View) selectRemoteFile(title string) (int, error) {
	var (
		strid string
		opts  []huh.Option[string]
	)
	for id, f := range v.gogh.Data.Filestore.Files {
		opts = append(opts, huh.NewOption[string](f.Filename, fmt.Sprint(id)))
	}
	if err := RunSelect(title, opts, &strid); err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(strid)
	if err != nil {
		return -1, fmt.Errorf("failed to get delete id")
	}
	return id, nil
}

func (v *View) deleteAction() {
	id, err := v.selectRemoteFile("delete")
	if err != nil {
		ErrorFunc(err)
		return
	}
	file := v.gogh.Data.Filestore.Files[id]
	if err := v.gogh.Delete(id); err != nil {
		ErrorFunc(err)
		return
	}
	fmt.Printf("✅ deleted %s\r\n", file.Filename)
}

func (v *View) downloadAction() {
	id, err := v.selectRemoteFile("download")
	if err != nil {
		ErrorFunc(err)
		return
	}
	file := v.gogh.Data.Filestore.Files[id]
	var download func() = func() {
		if err := v.gogh.Download(id); err != nil {
			ErrorFunc(err)
			return
		}
		fmt.Printf("✅ downloaded %s\r\n", file.Filename)
	}
	_ = spinner.New().Title(fmt.Sprintf("⏳ download %s", file.Filename)).Action(download).Run()
}

func (v *View) storeExecutor(in string) {
	args := strings.Fields(in)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case Upload:
		v.uploadAction()
	case Download:
		v.downloadAction()
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
