//go:build cgo
package gui

import (
	"context"
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/firasmosbehi/ssh-helper/internal/config"
	"github.com/firasmosbehi/ssh-helper/internal/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

func makeMCPView(w fyne.Window) fyne.CanvasObject {
	var servers map[string]config.MCPClientConfig
	var serverNames []string
	var selectedServer string
	var tools []mcpgo.Tool
	var selectedTool int = -1

	var serverList, toolList *widget.List
	var split *container.Split

	cfg, _ := config.NewManager()
	c, _ := cfg.Load()
	servers = c.MCPClients
	for name := range servers {
		serverNames = append(serverNames, name)
	}
	sort.Strings(serverNames)

	refreshTools := func() {
		selectedTool = -1
		if toolList != nil {
			toolList.UnselectAll()
			toolList.Refresh()
		}
	}

	serverList = widget.NewList(
		func() int { return len(serverNames) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(serverNames[id])
		},
	)
	serverList.OnSelected = func(id widget.ListItemID) {
		name := serverNames[id]
		selectedServer = name
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			t, err := mcp.ListTools(ctx, servers[name])
			if err != nil {
				fmt.Println("list tools error:", err)
				tools = nil
			} else {
				tools = t
			}
			refreshTools()
		}()
	}

	toolList = widget.NewList(
		func() int { return len(tools) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(fmt.Sprintf("%s - %s", tools[id].Name, tools[id].Description))
		},
	)
	toolList.OnSelected = func(id widget.ListItemID) { selectedTool = id }
	toolList.OnUnselected = func(id widget.ListItemID) { selectedTool = -1 }

	callBtn := widget.NewButton("Call Tool", func() {
		if selectedTool < 0 || selectedTool >= len(tools) || selectedServer == "" {
			return
		}
		t := tools[selectedTool]
		showMCPCallDialog(w, servers[selectedServer], t.Name)
	})

	right := container.NewBorder(callBtn, nil, nil, nil, toolList)
	split = container.NewHSplit(serverList, right)
	split.Offset = 0.3
	return split
}

func showMCPCallDialog(w fyne.Window, srv config.MCPClientConfig, toolName string) {
	argsEntry := widget.NewMultiLineEntry()
	argsEntry.SetText("{}")

	var d dialog.Dialog
	buttons := container.NewHBox(
		widget.NewButton("Cancel", func() { d.Hide() }),
		widget.NewButton("Call", func() {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				res, err := mcp.CallToolString(ctx, srv, toolName, argsEntry.Text)
				if err != nil {
					fmt.Println("call error:", err)
					return
				}
				result := widget.NewMultiLineEntry()
				result.SetText(res)
				result.Disable()
				dialog.ShowCustom("Result", "Close", result, w)
			}()
			d.Hide()
		}),
	)

	content := container.NewVBox(argsEntry, buttons)
	d = dialog.NewCustom(fmt.Sprintf("Call %s", toolName), "", content, w)
	d.Resize(fyne.NewSize(500, 300))
	d.Show()
}
