//go:build cgo
package gui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
)

func makeHostsView(w fyne.Window) fyne.CanvasObject {
	var hosts []core.Host
	var selectedID int = -1
	var list *widget.List

	cfgPath := os.ExpandEnv("$HOME/.ssh/config")

	refresh := func() {
		cfg, err := ssh.ParseConfig(cfgPath)
		if err == nil {
			hosts = ssh.HostsFromConfig(cfg)
		} else {
			hosts = nil
		}
		selectedID = -1
		if list != nil {
			list.UnselectAll()
			list.Refresh()
		}
	}
	refresh()

	list = widget.NewList(
		func() int { return len(hosts) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("template"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			h := hosts[id]
			label := item.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("%s  (%s@%s:%d)", h.Name, h.User, h.Hostname, h.Port))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selectedID = id
	}
	list.OnUnselected = func(id widget.ListItemID) {
		selectedID = -1
	}

	addBtn := widget.NewButton("Add", func() {
		showHostForm(w, cfgPath, core.Host{}, func() {
			refresh()
		})
	})

	editBtn := widget.NewButton("Edit", func() {
		if selectedID < 0 || selectedID >= len(hosts) {
			return
		}
		showHostForm(w, cfgPath, hosts[selectedID], func() {
			refresh()
		})
	})

	removeBtn := widget.NewButton("Remove", func() {
		if selectedID < 0 || selectedID >= len(hosts) {
			return
		}
		h := hosts[selectedID]
		d := dialog.NewConfirm("Remove Host", fmt.Sprintf("Remove %s?", h.Name), func(ok bool) {
			if ok {
				if err := ssh.RemoveHost(cfgPath, h.Name); err != nil {
					dialog.ShowError(err, w)
				} else {
					refresh()
				}
			}
		}, w)
		d.Show()
	})

	connectBtn := widget.NewButton("Connect", func() {
		if selectedID < 0 || selectedID >= len(hosts) {
			return
		}
		h := hosts[selectedID]
		dialog.ShowInformation("Connect", fmt.Sprintf("Connecting to %s (%s@%s:%d)...", h.Name, h.User, h.Hostname, h.Port), w)
	})

	buttons := container.NewHBox(addBtn, editBtn, removeBtn, connectBtn)
	return container.NewBorder(buttons, nil, nil, nil, list)
}
