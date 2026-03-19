package model

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestStylesDefined(t *testing.T) {
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"TitleStyle", TitleStyle},
		{"BodyStyle", BodyStyle},
		{"ChatBoxStyle", ChatBoxStyle},
		{"KanbanColumnStyle", KanbanColumnStyle},
		{"TicketCardStyle", TicketCardStyle},
		{"PopupStyle", PopupStyle},
		{"TabActiveStyle", TabActiveStyle},
		{"TabInactiveStyle", TabInactiveStyle},
	}
	for _, s := range styles {
		rendered := s.style.Render("test")
		if rendered == "" {
			t.Errorf("%s.Render returned empty string", s.name)
		}
	}
}

func TestColorsDefined(t *testing.T) {
	if ColorPrimary == "" {
		t.Error("ColorPrimary not defined")
	}
	if ColorSecondary == "" {
		t.Error("ColorSecondary not defined")
	}
	if ColorError == "" {
		t.Error("ColorError not defined")
	}
	if ColorSuccess == "" {
		t.Error("ColorSuccess not defined")
	}
}
