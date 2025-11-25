// kaf-mirror - A high-performance Kafka replication tool with AI-powered operational intelligence.
// Copyright (C) 2025 Scalytics
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package layouts

import (
	"kaf-mirror/cmd/mirror-cli/dashboard/core"
	"kaf-mirror/cmd/mirror-cli/dashboard/widgets"
	
	ui "github.com/gizak/termui/v3"
	ui_widgets "github.com/gizak/termui/v3/widgets"
)

type CategoryLayout struct {
	grid    *ui.Grid
	factory widgets.WidgetFactory
	context *core.NavigationContext
}

func NewCategoryLayout(factory widgets.WidgetFactory, context *core.NavigationContext) *CategoryLayout {
	grid := ui.NewGrid()
	
	return &CategoryLayout{
		grid:    grid,
		factory: factory,
		context: context,
	}
}

func (cl *CategoryLayout) Setup(termWidth, termHeight int) {
	cl.grid.SetRect(0, 0, termWidth, termHeight)
	
	if cl.factory == nil {
		return
	}
	
	listWidget := cl.factory.CreateListWidget()
	detailWidget := cl.factory.CreateDetailWidget("")
	
	cl.grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(0.4, listWidget),
			ui.NewCol(0.6, detailWidget),
		),
	)
}

func (cl *CategoryLayout) GetGrid() *ui.Grid {
	return cl.grid
}

func (cl *CategoryLayout) UpdateData(dataManager *core.DataManager) error {
	if cl.factory == nil {
		return nil
	}
	
	if err := cl.factory.UpdateListData(dataManager); err != nil {
		return err
	}
	
	if cl.context.ItemID != "" {
		if err := cl.factory.UpdateDetailData(dataManager, cl.context.ItemID); err != nil {
			return err
		}
		cl.factory.ResetDetailCursor()
	}
	
	return nil
}

func (cl *CategoryLayout) HandleScroll(direction core.Direction) {
	if cl.factory == nil {
		return
	}
	
	switch cl.context.PaneFocus {
	case core.ListPaneFocus:
		cl.factory.ScrollList(direction)
	case core.DetailPaneFocus:
		cl.factory.ScrollDetail(direction)
	}
}

func (cl *CategoryLayout) UpdatePaneFocus(focus core.PaneFocus) {
	if cl.factory == nil {
		return
	}
	
	listWidget := cl.factory.CreateListWidget()
	detailWidget := cl.factory.CreateDetailWidget("")
	
	cl.factory.SetDetailCursorVisible(focus == core.DetailPaneFocus)
	
	switch focus {
	case core.ListPaneFocus:
		if lw, ok := listWidget.(*ui_widgets.List); ok {
			lw.BorderStyle = ui.NewStyle(ui.ColorWhite)
		}
		if dw, ok := detailWidget.(*ui_widgets.List); ok {
			dw.BorderStyle = ui.NewStyle(ui.ColorYellow)
		}
	case core.DetailPaneFocus:
		if lw, ok := listWidget.(*ui_widgets.List); ok {
			lw.BorderStyle = ui.NewStyle(ui.ColorCyan)
		}
		if dw, ok := detailWidget.(*ui_widgets.List); ok {
			dw.BorderStyle = ui.NewStyle(ui.ColorWhite)
		}
	}
}

func (cl *CategoryLayout) HandleSelection() {
	if cl.factory == nil {
		return
	}
	
	selectedID := cl.factory.GetSelectedItemID()
	
	if selectedID != "" {
		cl.context.ItemID = selectedID
	}
}
