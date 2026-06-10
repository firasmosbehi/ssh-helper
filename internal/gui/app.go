package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Run launches the Fyne GUI application.
func Run() error {
	a := app.NewWithID("com.firasmosbehi.ssh-helper")
	w := a.NewWindow("ssh-helper")
	w.Resize(fyne.NewSize(900, 600))

	content := container.NewStack(widget.NewLabel("Select a section from the sidebar."))

	nav := container.NewVBox(
		widget.NewButtonWithIcon("Hosts", nil, func() {
			content.Objects = []fyne.CanvasObject{makeHostsView(w)}
			content.Refresh()
		}),
		widget.NewButtonWithIcon("Sync Jobs", nil, func() {
			content.Objects = []fyne.CanvasObject{makeJobsView(w)}
			content.Refresh()
		}),
		widget.NewButtonWithIcon("Keys", nil, func() {
			content.Objects = []fyne.CanvasObject{makeKeysView(w)}
			content.Refresh()
		}),
		widget.NewButtonWithIcon("MCP", nil, func() {
			content.Objects = []fyne.CanvasObject{makeMCPView(w)}
			content.Refresh()
		}),
	)

	split := container.NewHSplit(nav, content)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
	return nil
}
