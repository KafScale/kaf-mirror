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

package core

import (
	"time"
)

type NavigationState int

const (
	TopLevelState NavigationState = iota
	CategoryViewState
	DetailViewState
)

type Category int

const (
	ClustersCategory Category = iota
	JobsCategory
	InsightsCategory
	ComplianceCategory
)

type PaneFocus int

const (
	ListPaneFocus PaneFocus = iota
	DetailPaneFocus
)

type NavigationContext struct {
	State       NavigationState
	Category    Category
	ItemID      string
	Breadcrumbs []string
	LastUpdate  time.Time
	PaneFocus   PaneFocus
}

func NewNavigationContext() *NavigationContext {
	return &NavigationContext{
		State:       TopLevelState,
		Category:    ClustersCategory,
		ItemID:      "",
		Breadcrumbs: []string{"Dashboard"},
		LastUpdate:  time.Now(),
	}
}

func (nc *NavigationContext) NavigateToCategory(category Category) {
	nc.State = CategoryViewState
	nc.Category = category
	nc.ItemID = ""
	nc.PaneFocus = ListPaneFocus // Default to list pane
	
	categoryNames := map[Category]string{
		ClustersCategory:   "Clusters",
		JobsCategory:       "Jobs",
		InsightsCategory:   "Insights",
		ComplianceCategory: "Compliance",
	}
	
	nc.Breadcrumbs = []string{"Dashboard", categoryNames[category]}
	nc.LastUpdate = time.Now()
}

func (nc *NavigationContext) NavigateToDetail(itemID, itemName string) {
	nc.State = DetailViewState
	nc.ItemID = itemID
	
	if len(nc.Breadcrumbs) == 2 {
		nc.Breadcrumbs = append(nc.Breadcrumbs, itemName)
	} else if len(nc.Breadcrumbs) >= 3 {
		nc.Breadcrumbs[2] = itemName
	}
	
	nc.LastUpdate = time.Now()
}

func (nc *NavigationContext) NavigateBack() {
	switch nc.State {
	case DetailViewState:
		nc.State = CategoryViewState
		nc.ItemID = ""
		if len(nc.Breadcrumbs) > 2 {
			nc.Breadcrumbs = nc.Breadcrumbs[:2]
		}
	case CategoryViewState:
		nc.State = TopLevelState
		nc.Breadcrumbs = []string{"Dashboard"}
	case TopLevelState:
	}
	nc.LastUpdate = time.Now()
}

func (nc *NavigationContext) GetBreadcrumbText() string {
	text := ""
	for i, crumb := range nc.Breadcrumbs {
		if i > 0 {
			text += " > "
		}
		text += crumb
	}
	return text
}

func (nc *NavigationContext) IsAtTopLevel() bool {
	return nc.State == TopLevelState
}

func (nc *NavigationContext) IsInCategory() bool {
	return nc.State == CategoryViewState
}

func (nc *NavigationContext) IsInDetail() bool {
	return nc.State == DetailViewState
}
