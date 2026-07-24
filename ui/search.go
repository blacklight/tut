package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SearchView struct {
	tutView *TutView
	View    *tview.Flex
	Input   *tview.InputField
}

func NewSearchView(tv *TutView) *SearchView {
	input := NewInputField(tv.tut.Config)
	input.SetLabel("Search: ")
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			query := strings.TrimSpace(input.GetText())
			if len(query) > 0 {
				tv.SearchCommand(query)
			}
		}
		tv.FocusMainNoHistory()
		input.SetText("")
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(input, 1, 0, true).
		AddItem(nil, 0, 1, false)
	flex.SetBackgroundColor(tv.tut.Config.Style.Background)

	sv := &SearchView{
		tutView: tv,
		View:    flex,
		Input:   input,
	}
	return sv
}
