package view

import (
	"fmt"
	"gogh/pkg/fs"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
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
	}
	_ = spinner.New().Title(fmt.Sprintf("⏳ upload %s", file)).Action(upload).Run()
	fmt.Printf("👌 uploaded %s\n", file)
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
		for _, v := range v.gogh.Data.Filestore.Files {
			fmt.Println(v.Filename)
			for _, p := range v.Pieces {
				fmt.Println(p.URL)
			}
		}
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
