package style

import (
	"github.com/charmbracelet/bubbles/v2/list"
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

var (
	listTitleBarStyle = lipgloss.NewStyle()

	listTitleStyle = TextStyle.
			Foreground(LightBlue).
			Bold(true)


	listSpinnerStyle                     = lipgloss.NewStyle().Foreground(LightBlue)
	listFilterPromptStyle                = lipgloss.NewStyle().Foreground(LightGreen)
	listFilterCursorStyle                = lipgloss.NewStyle().Foreground(LightGreen)
	listDefaultFilterCharacterMatchStyle = lipgloss.NewStyle().Foreground(LightGreen)


	listStatusBarStyle = lipgloss.NewStyle().
				Foreground(DarkBlue).
				Margin(0, 0, 1, 0)

	listStatusEmptyStyle           = lipgloss.NewStyle().Foreground(LightGreen)
	listStatusBarActiveFilterStyle = lipgloss.NewStyle().Foreground(LightGreen)
	listStatusBarFilterCountStyle  = lipgloss.NewStyle().Foreground(LightGreen)
	listNoItemsStyle               = lipgloss.NewStyle().Foreground(LightGreen)
	listPaginationStyleStyle       = lipgloss.NewStyle().Foreground(LightGreen)
	listHelpStyleStyle             = lipgloss.NewStyle().Foreground(LightGreen)
	listActivePaginationDotStyle   = lipgloss.NewStyle().Foreground(LightGreen)
	listInactivePaginationDotStyle = lipgloss.NewStyle().Foreground(LightGreen)
	listArabicPaginationStyle      = lipgloss.NewStyle().Foreground(LightGreen)
	listDividerDotStyle            = lipgloss.NewStyle().Foreground(LightGreen)


	listDelegateNormalTitleStyle = TextStyle
	listDelegateNormalDescStyle  = TextStyle

	listDelegateSelectedlTitleStyle = lipgloss.NewStyle().
					Foreground(LightGreen).
					Bold(true)

	listDelegateSelectedDescStyle = lipgloss.NewStyle().
		Foreground(LightGreen)

	listDelegateDimmedTitleStyle = lipgloss.NewStyle().Foreground(LightBlack)
	listDelegateDimmedDescStyle  = lipgloss.NewStyle().Foreground(LightBlack)

	listDelegateFilterMatchStyle = lipgloss.NewStyle().Foreground(LightGreen)

	// dimmed styles

	listTitleStyleDimmed = TextStyle.
		Foreground(LightBlueDimmed).
		Bold(true)

	listSpinnerStyleDimmed                     = lipgloss.NewStyle().Foreground(LightBlueDimmed)
	listFilterPromptStyleDimmed                = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listFilterCursorStyleDimmed                = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listDefaultFilterCharacterMatchStyleDimmed = lipgloss.NewStyle().Foreground(LightGreenDimmed)

	listStatusBarStyleDimmed = lipgloss.NewStyle().
		Foreground(DarkBlueDimmed).
		Margin(0, 0, 1, 0)

	listStatusEmptyStyleDimmed           = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listStatusBarActiveFilterStyleDimmed = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listStatusBarFilterCountStyleDimmed  = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listNoItemsStyleDimmed               = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listPaginationStyleStyleDimmed       = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listHelpStyleStyleDimmed             = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listActivePaginationDotStyleDimmed   = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listInactivePaginationDotStyleDimmed = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listArabicPaginationStyleDimmed      = lipgloss.NewStyle().Foreground(LightGreenDimmed)
	listDividerDotStyleDimmed            = lipgloss.NewStyle().Foreground(LightGreenDimmed)


	listDelegateNormalTitleStyleDimmed = TextStyleDimmed
	listDelegateNormalDescStyleDimmed  = TextStyleDimmed

	listDelegateSelectedlTitleStyleDimmed = lipgloss.NewStyle().
		Foreground(LightGreenDimmed).
		Bold(true)

	listDelegateSelectedDescStyleDimmed = lipgloss.NewStyle().
		Foreground(LightGreenDimmed)

	listDelegateFilterMatchStyleDimmed = lipgloss.NewStyle().Foreground(LightGreen)

)

func DefaultListStyles() list.Styles {
	s := list.DefaultStyles(true)

	s.TitleBar = listTitleBarStyle
	s.Title = listTitleStyle
	s.Spinner = listSpinnerStyle

	s.DefaultFilterCharacterMatch = listDefaultFilterCharacterMatchStyle

	s.StatusBar = listStatusBarStyle
	s.StatusEmpty = listStatusEmptyStyle
	s.StatusBarActiveFilter = listStatusBarActiveFilterStyle
	s.StatusBarFilterCount = listStatusBarFilterCountStyle

	s.NoItems = listNoItemsStyle

	s.PaginationStyle = listPaginationStyleStyle
	s.HelpStyle = listHelpStyleStyle

	s.ActivePaginationDot = listActivePaginationDotStyle
	s.InactivePaginationDot = listInactivePaginationDotStyle
	s.ArabicPagination = listArabicPaginationStyle
	s.DividerDot = listDividerDotStyle

	return s
}

func DefaultItemStyles() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles(true)

	s.NormalTitle = listDelegateNormalTitleStyle
	s.NormalDesc = listDelegateNormalDescStyle

	s.SelectedTitle = listDelegateSelectedlTitleStyle
	s.SelectedDesc = listDelegateSelectedDescStyle

	s.DimmedTitle = listDelegateDimmedTitleStyle
	s.DimmedDesc = listDelegateDimmedDescStyle

	s.FilterMatch = listDelegateFilterMatchStyle

	return s
}

func DefaultListStylesDimmed() list.Styles {
	s := list.DefaultStyles(true)

	s.TitleBar = listTitleBarStyle
	s.Title = listTitleStyleDimmed
	s.Spinner = listSpinnerStyleDimmed

	s.DefaultFilterCharacterMatch = listDefaultFilterCharacterMatchStyleDimmed

	s.StatusBar = listStatusBarStyleDimmed
	s.StatusEmpty = listStatusEmptyStyleDimmed
	s.StatusBarActiveFilter = listStatusBarActiveFilterStyleDimmed
	s.StatusBarFilterCount = listStatusBarFilterCountStyleDimmed

	s.NoItems = listNoItemsStyleDimmed

	s.PaginationStyle = listPaginationStyleStyleDimmed
	s.HelpStyle = listHelpStyleStyleDimmed

	s.ActivePaginationDot = listActivePaginationDotStyleDimmed
	s.InactivePaginationDot = listInactivePaginationDotStyleDimmed
	s.ArabicPagination = listArabicPaginationStyleDimmed
	s.DividerDot = listDividerDotStyleDimmed

	return s
}


func DefaultItemStylesDimmed() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles(true)

	s.NormalTitle = listDelegateNormalTitleStyleDimmed
	s.NormalDesc = listDelegateNormalDescStyleDimmed

	s.SelectedTitle = listDelegateSelectedlTitleStyleDimmed
	s.SelectedDesc = listDelegateSelectedDescStyleDimmed

	s.DimmedTitle = listDelegateDimmedTitleStyle
	s.DimmedDesc = listDelegateDimmedDescStyle

	s.FilterMatch = listDelegateFilterMatchStyleDimmed

	return s
}
