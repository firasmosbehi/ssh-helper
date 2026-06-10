package gui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
)

func makeKeysView(w fyne.Window) fyne.CanvasObject {
	var keys []ssh.KeyInfo
	var selectedID int = -1
	var list *widget.List

	sshDir := os.ExpandEnv("$HOME/.ssh")

	refresh := func() {
		k, err := ssh.ListKeys(sshDir)
		if err == nil {
			keys = k
		} else {
			keys = nil
		}
		selectedID = -1
		if list != nil {
			list.UnselectAll()
			list.Refresh()
		}
	}
	refresh()

	list = widget.NewList(
		func() int { return len(keys) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("template"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			k := keys[id]
			label := item.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("%s  (%s)", k.Name, k.Fingerprint))
		},
	)

	list.OnSelected = func(id widget.ListItemID) { selectedID = id }
	list.OnUnselected = func(id widget.ListItemID) { selectedID = -1 }

	genBtn := widget.NewButton("Generate", func() {
		showKeyGenForm(w, sshDir, refresh)
	})

	removeBtn := widget.NewButton("Remove", func() {
		if selectedID < 0 || selectedID >= len(keys) {
			return
		}
		k := keys[selectedID]
		dialog.NewConfirm("Remove Key", fmt.Sprintf("Remove %s?", k.Name), func(ok bool) {
			if ok {
				if err := ssh.RemoveKey(sshDir + "/" + ssh.SanitizeKeyName(k.Name)); err != nil {
					dialog.ShowError(err, w)
				} else {
					refresh()
				}
			}
		}, w).Show()
	})

	buttons := container.NewHBox(genBtn, removeBtn)
	return container.NewBorder(buttons, nil, nil, nil, list)
}

func showKeyGenForm(w fyne.Window, sshDir string, onSave func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("id_rsa")

	typeSelect := widget.NewSelect([]string{"ed25519", "rsa"}, func(s string) {})
	typeSelect.SetSelected("ed25519")

	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Type", typeSelect),
	)

	var d dialog.Dialog
	buttons := container.NewHBox(
		widget.NewButton("Cancel", func() { d.Hide() }),
		widget.NewButton("Generate", func() {
			name := nameEntry.Text
			if name == "" {
				lbl := widget.NewLabel("Name is required")
				lbl.Importance = widget.DangerImportance
				dialog.ShowCustom("Validation", "OK", lbl, w)
				return
			}
			if err := ssh.GenerateKey(sshDir, name, typeSelect.Selected); err != nil {
				dialog.ShowError(err, w)
				return
			}
			d.Hide()
			onSave()
		}),
	)

	content := container.NewVBox(form, buttons)
	d = dialog.NewCustom("Generate Key", "", content, w)
	d.Resize(fyne.NewSize(350, 200))
	d.Show()
}
