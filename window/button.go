package window

import (
	"fyne.io/fyne/v2/widget"
)

type customButton struct {
	*widget.Button
	btnTag int
}

func newCustomButton(text string) *customButton {
	btn := widget.NewButton(text, nil)
	btn.Alignment = widget.ButtonAlignCenter
	btn.Importance = widget.LowImportance
	return &customButton{
		Button: btn,
	}
}

func (btn *customButton) setTag(tag int) {
	btn.btnTag = tag
}

func (btn *customButton) getTag() int {
	return btn.btnTag
}
