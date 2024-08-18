package options

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func PromptOptions(prefix string, exit string) []prompt.Option {
	options := []prompt.Option{
		prompt.OptionPrefix(prefix),
		prompt.OptionTitle(prefix),
		prompt.OptionPrefixTextColor(0),
		prompt.OptionPrefixTextColor(prompt.DefaultColor),
		prompt.OptionPrefixBackgroundColor(prompt.DefaultColor),
		prompt.OptionInputTextColor(prompt.DefaultColor),
		prompt.OptionInputBGColor(prompt.DefaultColor),
		prompt.OptionPreviewSuggestionTextColor(prompt.DefaultColor),
		prompt.OptionPreviewSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionSuggestionTextColor(prompt.DefaultColor),
		prompt.OptionSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionSelectedSuggestionTextColor(prompt.Green),
		prompt.OptionSelectedSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionDescriptionTextColor(prompt.DefaultColor),
		prompt.OptionDescriptionBGColor(prompt.DefaultColor),
		prompt.OptionSelectedDescriptionTextColor(prompt.Green),
		prompt.OptionSelectedDescriptionBGColor(prompt.DefaultColor),
		prompt.OptionScrollbarThumbColor(prompt.DefaultColor),
		prompt.OptionScrollbarBGColor(prompt.DefaultColor),
		prompt.OptionMaxSuggestion(3),
		prompt.OptionSetExitCheckerOnInput(prompt.ExitChecker(func(in string, breakline bool) bool {
			return strings.TrimSpace(in) == exit
		})),
	}
	return options
}
