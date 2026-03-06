package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	left := m.theme.Pane.Width(m.width/3 - 1).Height(m.height - 4).Render(m.list.View())
	rightContent := m.actionStateView() + "\n\n" + m.detailText
	if m.pendingConfirm != nil {
		rightContent = rightContent + "\n\n" +
			"Confirm: " + m.pendingConfirm.Message + "\n" +
			"Press y/enter to confirm, n/esc to cancel"
	}
	if strings.TrimSpace(m.errText) != "" {
		rightContent = "Status: " + m.statusText + "\n\n" + rightContent
	}
	right := m.theme.Pane.Width(m.width - m.width/3 - 3).Height(m.height - 4).Render(rightContent)
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	parts := []string{}
	if m.cfg.UI.ShowHeader {
		parts = append(parts, m.headerView())
	}
	parts = append(parts, body)
	if m.cfg.UI.ShowFooter {
		parts = append(parts, m.footerView())
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) actionStateView() string {
	if m.cfg.ReadOnly {
		return "Actions: readonly mode active (x delete, p patch, k cancel disabled)"
	}

	parts := []string{}
	if m.permissions.DeleteAllowed {
		parts = append(parts, "x=enabled")
	} else {
		reason := m.permissions.DeleteReason
		if strings.TrimSpace(reason) == "" {
			reason = "RBAC"
		}
		parts = append(parts, "x=disabled("+reason+")")
	}

	if m.permissions.PatchAllowed {
		parts = append(parts, "p=enabled")
	} else {
		reason := m.permissions.PatchReason
		if strings.TrimSpace(reason) == "" {
			reason = "RBAC"
		}
		parts = append(parts, "p=disabled("+reason+")")
	}

	if m.permissions.CancelAllowed {
		parts = append(parts, "k=enabled")
	} else {
		reason := m.permissions.CancelReason
		if strings.TrimSpace(reason) == "" {
			reason = "RBAC or Tekton API unavailable"
		}
		parts = append(parts, "k=disabled("+reason+")")
	}

	return "Actions: " + strings.Join(parts, "  ")
}
