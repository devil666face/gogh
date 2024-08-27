package view

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

func (v *View) exportAction() {
	id, err := v.selectRemoteFile("export")
	if err != nil {
		ErrorFunc(err)
		return
	}
	file := v.gogh.Data.Filestore.Files[id]
	if err := v.gogh.SaveFile(&file); err != nil {
		ErrorFunc(err)
		return
	}
	fmt.Printf("💾 file %s saved\r\n", file.Filename)
}

func (v *View) importAction() {
	filename, err := v.selectLocalFile("import .gogh file", ".gogh")
	if err != nil {
		ErrorFunc(err)
		return
	}
	file, err := v.gogh.LoadFile(filename)
	v.gogh.Data.Filestore.Files = append(v.gogh.Data.Filestore.Files, *file)
	v.gogh.Data.Filestore.Sort()
	if err := v.gogh.SaveData(); err != nil {
		ErrorFunc(err)
		return
	}
	fmt.Printf("📥 file %s imported\r\n", file.Filename)
}

func (v *View) shareExecutor(in string) {
	args := strings.Fields(in)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case Export:
		v.exportAction()
	case Import:
		v.importAction()
	default:
		NotFoundFunc()
	}
}

func (v *View) shareCompleter(d prompt.Document) []prompt.Suggest {
	var complete = []prompt.Suggest{
		{Text: Export, Description: "export file for share from remote store"},
		{Text: Import, Description: "import shared file in local store"},
	}
	// remove second postition complete
	for _, c := range complete {
		if HasPrefix(d, c.Text) {
			complete = []prompt.Suggest{}
		}
	}

	return prompt.FilterContains(complete, d.GetWordBeforeCursor(), true)
}
