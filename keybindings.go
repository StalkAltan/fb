package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

func (a *FileBundlerApp) setupKeyBindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focusedPrimitive := a.app.GetFocus()

		// Check if we're in a modal context (either Modal or its Button)
		_, isModal := focusedPrimitive.(*tview.Modal)
		_, isButton := focusedPrimitive.(*tview.Button)
		if isModal || isButton {
			return event
		}

		if _, ok := a.app.GetFocus().(*tview.Modal); ok {
			return event
		}

		// Handle the help toggle separately, before any other checks
		if event.Rune() == '?' {
			a.toggleHelpPanel()
			return nil
		}

		node := a.ui.fileTree.GetCurrentNode()
		if node == nil {
			return event
		}

		// Handle global keys
		switch event.Key() {
		case tcell.KeyTab:
			if a.app.GetFocus() == a.ui.fileTree {
				a.app.SetFocus(a.ui.selectionList)
			} else {
				a.app.SetFocus(a.ui.fileTree)
			}
			return nil
		case tcell.KeyEnter:
			path := node.GetReference().(string)
			info, err := os.Stat(path)
			if err != nil {
				return event
			}

			if info.IsDir() {
				if node.IsExpanded() {
					// Collapse directory
					node.SetExpanded(false)
					a.ui.statusBar.SetText("Collapsed: " + path)
				} else {
					// Populate on first expansion
					if len(node.GetChildren()) == 0 {
						err := a.populateNode(node)
						if err != nil {
							a.ui.statusBar.SetText("Error loading directory: " + err.Error())
							return nil
						}
					}
					node.SetExpanded(true)
					a.ui.statusBar.SetText("Expanded: " + path)
				}
				return nil
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ': // Space key
				path := node.GetReference().(string)
				info, err := os.Stat(path)
				if err != nil {
					return nil
				}

				if info.IsDir() {
					// Check if any files in the directory are already selected
					isAnySelected := false
					var checkSelection func(*tview.TreeNode) bool
					checkSelection = func(n *tview.TreeNode) bool {
						for _, child := range n.GetChildren() {
							childPath := child.GetReference().(string)
							childInfo, err := os.Stat(childPath)
							if err != nil {
								continue
							}

							if childInfo.IsDir() {
								if checkSelection(child) {
									return true
								}
							} else {
								if a.selection.IsSelected(childPath) {
									return true
								}
							}
						}
						return false
					}
					isAnySelected = checkSelection(node)

					// Function to recursively select all files in a directory
					var selectAllInDir func(*tview.TreeNode, bool)
					selectAllInDir = func(n *tview.TreeNode, selecting bool) {
						// Ensure the directory is populated
						if len(n.GetChildren()) == 0 {
							err := a.populateNode(n)
							if err != nil {
								return
							}
						}

						// Process all children
						for _, child := range n.GetChildren() {
							childPath := child.GetReference().(string)
							childInfo, err := os.Stat(childPath)
							if err != nil {
								continue
							}

							if childInfo.IsDir() {
								selectAllInDir(child, selecting)
							} else {
								if selecting {
									a.selection.Add(childPath)
									child.SetColor(tcell.ColorGreen.TrueColor())
								} else {
									a.selection.Remove(childPath)
									child.SetColor(tcell.ColorWhite.TrueColor())
								}
							}
						}
					}

					// Start recursive selection/deselection
					selectAllInDir(node, !isAnySelected)

					// Update status bar
					if !isAnySelected {
						a.ui.statusBar.SetText("Selected all files in: " + path)
					} else {
						a.ui.statusBar.SetText("Deselected all files in: " + path)
					}
				} else {
					// Handle single file selection
					a.selection.Toggle(path)
					if a.selection.IsSelected(path) {
						node.SetColor(tcell.ColorGreen.TrueColor())
						a.ui.statusBar.SetText("Selected: " + path)
					} else {
						node.SetColor(tcell.ColorWhite.TrueColor())
						a.ui.statusBar.SetText("Deselected: " + path)
					}
				}
				// Update selection list
				a.updateSelectionList()
				return nil
			case 'q':
				a.app.Stop()
				return nil
			case 'x':
				if len(a.selection.GetSelectedPaths()) > 0 {
					a.showBundleConfirmation()
				} else {
					a.showMessage("No files selected. Select files before bundling.", true)
				}
				return nil
			case 'k':
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			case 'j':
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 'h':
				return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
			case 'l':
				return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
			}
		}

		return event
	})

	// Add change handler to update status bar
	a.ui.fileTree.SetChangedFunc(func(node *tview.TreeNode) {
		if node != nil {
			path := node.GetReference().(string)
			a.ui.statusBar.SetText("Current: " + path)
		}
	})
}
