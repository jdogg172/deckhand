package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Header   lipgloss.Style
	Footer   lipgloss.Style
	Pane     lipgloss.Style
	Selected lipgloss.Style
	Healthy  lipgloss.Style
	Warning  lipgloss.Style
	Error    lipgloss.Style
	Muted    lipgloss.Style
	Title    lipgloss.Style
}

func DefaultTheme(noColor bool) Theme {
	if noColor {
		base := lipgloss.NewStyle()
		return Theme{
			Header:   base.Bold(true),
			Footer:   base,
			Pane:     base,
			Selected: base.Bold(true),
			Healthy:  base,
			Warning:  base,
			Error:    base,
			Muted:    base,
			Title:    base.Bold(true),
		}
	}

	return Theme{
		Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("62")).Padding(0, 1),
		Footer:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Pane:     lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1),
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("63")),
		Healthy:  lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true),
		Warning:  lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
		Error:    lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
		Muted:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		Title:    lipgloss.NewStyle().Bold(true).Underline(true),
	}
}
