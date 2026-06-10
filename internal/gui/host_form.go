//go:build cgo
package gui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
)

func showHostForm(w fyne.Window, cfgPath string, host core.Host, onSave func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(host.Name)
	hostEntry := widget.NewEntry()
	hostEntry.SetText(host.Hostname)
	userEntry := widget.NewEntry()
	userEntry.SetText(host.User)
	portEntry := widget.NewEntry()
	portEntry.SetText(strconv.Itoa(host.Port))
	if host.Port == 0 {
		portEntry.SetText("22")
	}
	keyEntry := widget.NewEntry()
	keyEntry.SetText(host.IdentityFile)

	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Hostname", hostEntry),
		widget.NewFormItem("User", userEntry),
		widget.NewFormItem("Port", portEntry),
		widget.NewFormItem("Identity File", keyEntry),
	)

	var d dialog.Dialog

	buttons := container.NewHBox(
		widget.NewButton("Cancel", func() {
			d.Hide()
		}),
		widget.NewButton("Save", func() {
			port, _ := strconv.Atoi(portEntry.Text)
			if port == 0 {
				port = 22
			}
			h := core.Host{
				Name:         nameEntry.Text,
				Hostname:     hostEntry.Text,
				User:         userEntry.Text,
				Port:         port,
				IdentityFile: keyEntry.Text,
			}
			if h.Name == "" || h.Hostname == "" {
				lbl := widget.NewLabel("Name and Hostname are required")
				lbl.Importance = widget.DangerImportance
				dialog.ShowCustom("Validation", "OK", lbl, w)
				return
			}
			if err := ssh.AddHost(cfgPath, h); err != nil {
				dialog.ShowError(err, w)
				return
			}
			d.Hide()
			onSave()
		}),
	)

	content := container.NewVBox(form, buttons)
	d = dialog.NewCustom("Host", "", content, w)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}
