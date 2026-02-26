package screens

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestLoginModel_InitialView(t *testing.T) {
	m := NewLoginModel()
	m.width = 80
	m.height = 24
	view := m.View()
	assert.Contains(t, view, "WHOOP CLI")
	assert.Contains(t, view, "log in")
}

func TestLoginModel_Loading(t *testing.T) {
	m := NewLoginModel()
	m.width = 80
	m.height = 24

	m, _ = m.Update(LoginMsg{})
	assert.True(t, m.loading)

	view := m.View()
	assert.Contains(t, view, "Waiting for browser")
}

func TestLoginModel_Error(t *testing.T) {
	m := NewLoginModel()
	m.width = 80
	m.height = 24

	m, _ = m.Update(LoginErrorMsg{Err: errors.New("test error")})
	view := m.View()
	assert.Contains(t, view, "Login failed")
	assert.Contains(t, view, "test error")
}

func TestLoginModel_WindowResize(t *testing.T) {
	m := NewLoginModel()
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
}
