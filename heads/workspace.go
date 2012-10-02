package heads

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/xgbutil/xrect"

	"github.com/BurntSushi/wingo/logger"
	"github.com/BurntSushi/wingo/workspace"
)

func (hds *Heads) Workspaces() []*workspace.Workspace {
	return hds.workspaces.Wrks
}

// ActivateWorkspace will "focus" or "activate" the workspace provided.
// This only works when "wk" is visible.
// To activate a hidden workspace, please use SwitchWorkspaces.
func (hds *Heads) ActivateWorkspace(wk *workspace.Workspace) {
	wkvi := hds.visibleIndex(wk)
	if wkvi > -1 {
		hds.active = wkvi
	}
}

func (hds *Heads) SwitchWorkspaces(wk1, wk2 *workspace.Workspace) {
	v1, v2 := hds.visibleIndex(wk1), hds.visibleIndex(wk2)
	switch {
	case v1 > -1 && v2 > -1:
		wk1.Hide()
		wk2.Hide()
		hds.visibles[v1], hds.visibles[v2] = hds.visibles[v2], hds.visibles[v1]
		wk1.Place()
		wk2.Place()
		wk1.Show()
		wk2.Show()
	case v1 > -1 && v2 == -1:
		wk1.Hide()
		hds.visibles[v1] = wk2
		wk2.Show()
	case v1 == -1 && v2 > -1:
		wk2.Hide()
		hds.visibles[v2] = wk1
		wk1.Show()
	case v1 == -1 && v2 == -1:
		// Meaningless
	default:
		panic("unreachable")
	}
}

func (hds *Heads) NewWorkspace(name string) *workspace.Workspace {
	return hds.workspaces.NewWorkspace(name)
}

func (hds *Heads) AddWorkspace(wk *workspace.Workspace) {
	hds.workspaces.Add(wk)
}

func (hds *Heads) RemoveWorkspace(wk *workspace.Workspace) {
	// Don't allow it if this would result in fewer workspaces than there
	// are active physical heads.
	if len(hds.geom) == len(hds.workspaces.Wrks) {
		return
	}

	// There's a bit of complexity in choosing where to move the clients to.
	// Namely, if we're removing a hidden workspace, it's a simple matter of
	// moving the clients. However, if we're removing a visible workspace,
	// we have to make sure to make another workspace that is hidden take
	// its place. (Such a workspace is guaranteed to exist because we have at
	// least one more workspace than there are active physical heads.)
	if !wk.IsVisible() {
		moveClientsTo := hds.workspaces.Wrks[len(hds.workspaces.Wrks)-1]
		if moveClientsTo == wk {
			moveClientsTo = hds.workspaces.Wrks[len(hds.workspaces.Wrks)-2]
		}
		wk.RemoveAllAndAdd(moveClientsTo)
	} else {
		// Find the last-most hidden workspace that is not itself.
		for i := len(hds.workspaces.Wrks) - 1; i >= 0; i-- {
			work := hds.workspaces.Wrks[i]
			if work != wk && !work.IsVisible() {
				hds.SwitchWorkspaces(wk, work)
				wk.RemoveAllAndAdd(work)
				break
			}
		}
	}
	hds.workspaces.Remove(wk)
}

func (hds *Heads) ActiveWorkspace() *workspace.Workspace {
	return hds.visibles[hds.active]
}

func (hds *Heads) VisibleWorkspaces() []*workspace.Workspace {
	return hds.visibles
}

// WithVisibleWorkspace takes a head number and a closure and executes
// the closure safely with the workspace corresponding to head number i.
//
// This approach is necessary for safety, since the user can send commands
// with arbitrary head numbers. We need to make sure we don't crash if we
// get an invalid head number.
func (hds *Heads) WithVisibleWorkspace(i int, f func(w *workspace.Workspace)) {
	if i < 0 || i >= len(hds.visibles) {
		headNums := make([]string, len(hds.visibles))
		for j := range headNums {
			headNums[j] = fmt.Sprintf("%d", j)
		}
		logger.Warning.Printf("Head index %d is not valid. "+
			"Valid heads are: [%s].", i, strings.Join(headNums, ", "))
		return
	}
	f(hds.visibles[i])
}

func (hds *Heads) FindMostOverlap(needle xrect.Rect) *workspace.Workspace {
	haystack := make([]xrect.Rect, len(hds.geom))
	for i := range haystack {
		haystack[i] = hds.geom[i]
	}

	index := xrect.LargestOverlap(needle, haystack)
	if index == -1 {
		return nil
	}
	return hds.visibles[index]
}

func (hds *Heads) IsActive(wrk *workspace.Workspace) bool {
	return hds.visibles[hds.active] == wrk
}

func (hds *Heads) Geom(wrk *workspace.Workspace) xrect.Rect {
	vi := hds.visibleIndex(wrk)
	if vi >= 0 {
		return hds.workarea[vi]
	}
	return nil
}

func (hds *Heads) visibleIndex(wk *workspace.Workspace) int {
	for i, vwk := range hds.visibles {
		if vwk == wk {
			return i
		}
	}
	return -1
}
