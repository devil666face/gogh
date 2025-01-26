package view

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (v *View) showSettingsAction() {
	var (
		sb strings.Builder
		w  = tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
	)
	fmt.Fprintf(w, "%s\t%s", "Key", "Value")
	fmt.Fprintf(w, "\n%s\t%s", "compress", boolFormat(v.gogh.Data.Settings.Compress))
	fmt.Fprintf(w, "\n%s\t%s", "token", v.gogh.Data.Settings.SessionToken)
	fmt.Fprintf(w, "\n%s\t%s", "device", v.gogh.Data.Settings.DeviceID)
	w.Flush()
	fmt.Println(
		lipgloss.NewStyle().
			Padding(0, 1).
			Render(sb.String()),
	)
}

var strValidator = func(in string) error {
	if in == "" {
		return fmt.Errorf("value is empty")
	}
	return nil
}

func (v *View) settingsExecutor(in string) {
	args := strings.Fields(in)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case Show:
		v.showSettingsAction()
	case Set:
		if len(args) == 1 {
			ErrorFunc(fmt.Errorf("you must set settings key"))
			return
		}
		switch args[1] {
		case Token:
			var (
				token string
			)
			if err := huh.NewInput().
				Title("<user_session>").
				Prompt(">> ").
				Validate(strValidator).
				Value(&token).
				Placeholder(v.gogh.Data.Settings.SessionToken).
				Run(); err != nil {
				ErrorFunc(err)
				return
			}
			if err := v.gogh.SetToken(token); err != nil {
				ErrorFunc(err)
			}
		case Device:
			var (
				deviceID string
			)
			if err := huh.NewInput().
				Title("<_device_id>").
				Prompt(">> ").
				Validate(strValidator).
				Value(&deviceID).
				Placeholder(v.gogh.Data.Settings.DeviceID).
				Run(); err != nil {
				ErrorFunc(err)
				return
			}
			if err := v.gogh.SetDeviceID(deviceID); err != nil {
				ErrorFunc(err)
			}
		case Compress:
			var confirm bool
			if err := huh.NewConfirm().
				Title("enable auto-compress").
				Affirmative("yes").
				Negative("no").
				Value(&confirm).
				Run(); err != nil {
				ErrorFunc(err)
				return
			}
			v.gogh.Data.Settings.Compress = confirm
			if err := v.gogh.SaveData(); err != nil {
				ErrorFunc(err)
			}
		}

	default:
		NotFoundFunc()
	}
}

func (v *View) settingsCompleter(d prompt.Document) []prompt.Suggest {
	var complete = []prompt.Suggest{
		{Text: Show, Description: "show all settings"},
		{Text: Set, Description: "set settings"},
	}
	// remove second postition complete
	for _, c := range complete {
		if HasPrefix(d, c.Text) {
			complete = []prompt.Suggest{}
		}
	}
	if HasPrefix(d, Set) {
		complete = []prompt.Suggest{
			{Text: Compress, Description: "set default compress or no"},
			{Text: Token, Description: "set github <user_session>"},
			{Text: Device, Description: "set github <_device_id>"},
		}
	}

	return prompt.FilterContains(complete, d.GetWordBeforeCursor(), true)
}
