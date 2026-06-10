package gui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/rsync"
	"github.com/firasmosbehi/ssh-helper/internal/store"
	"github.com/google/uuid"
)

func makeJobsView(w fyne.Window) fyne.CanvasObject {
	var jobs []core.SyncJob
	var selectedID int = -1
	var list *widget.List

	dir := os.ExpandEnv("$HOME/.config/ssh-helper")
	s := store.NewJSONStore(dir)

	refresh := func() {
		j, err := s.ListSyncJobs()
		if err == nil {
			jobs = j
		} else {
			jobs = nil
		}
		selectedID = -1
		if list != nil {
			list.UnselectAll()
			list.Refresh()
		}
	}
	refresh()

	list = widget.NewList(
		func() int { return len(jobs) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("template"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			j := jobs[id]
			label := item.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("%s  %s -> %s", j.Name, j.Source, j.Dest))
		},
	)

	list.OnSelected = func(id widget.ListItemID) { selectedID = id }
	list.OnUnselected = func(id widget.ListItemID) { selectedID = -1 }

	addBtn := widget.NewButton("Add", func() {
		showJobForm(w, s, core.SyncJob{}, refresh)
	})

	removeBtn := widget.NewButton("Remove", func() {
		if selectedID < 0 || selectedID >= len(jobs) {
			return
		}
		j := jobs[selectedID]
		dialog.NewConfirm("Remove Job", fmt.Sprintf("Remove %s?", j.Name), func(ok bool) {
			if ok {
				if err := s.DeleteSyncJob(j.ID); err != nil {
					dialog.ShowError(err, w)
				} else {
					refresh()
				}
			}
		}, w).Show()
	})

	runBtn := widget.NewButton("Run", func() {
		if selectedID < 0 || selectedID >= len(jobs) {
			return
		}
		j := jobs[selectedID]
		go func() {
			runner := rsync.Runner{Job: j}
			_ = runner.Run(nil, rsync.RunOptions{})
		}()
		dialog.ShowInformation("Run", fmt.Sprintf("Started job %s", j.Name), w)
	})

	buttons := container.NewHBox(addBtn, removeBtn, runBtn)
	return container.NewBorder(buttons, nil, nil, nil, list)
}

func showJobForm(w fyne.Window, s *store.JSONStore, job core.SyncJob, onSave func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(job.Name)
	srcEntry := widget.NewEntry()
	srcEntry.SetText(job.Source)
	dstEntry := widget.NewEntry()
	dstEntry.SetText(job.Dest)

	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Source", srcEntry),
		widget.NewFormItem("Destination", dstEntry),
	)

	var d dialog.Dialog
	buttons := container.NewHBox(
		widget.NewButton("Cancel", func() { d.Hide() }),
		widget.NewButton("Save", func() {
			if nameEntry.Text == "" || srcEntry.Text == "" || dstEntry.Text == "" {
				lbl := widget.NewLabel("All fields are required")
				lbl.Importance = widget.DangerImportance
				dialog.ShowCustom("Validation", "OK", lbl, w)
				return
			}
			j := core.SyncJob{
				ID:     job.ID,
				Name:   nameEntry.Text,
				Source: srcEntry.Text,
				Dest:   dstEntry.Text,
			}
			if j.ID == "" {
				j.ID = uuid.New().String()
			}
			if err := s.SaveSyncJob(j); err != nil {
				dialog.ShowError(err, w)
				return
			}
			d.Hide()
			onSave()
		}),
	)

	content := container.NewVBox(form, buttons)
	d = dialog.NewCustom("Sync Job", "", content, w)
	d.Resize(fyne.NewSize(400, 250))
	d.Show()
}
