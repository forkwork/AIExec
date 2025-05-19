package plugin

import (
	"fmt"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"aiexec-plugin/pkg/entities/manifest_entities"
)

type versionRequire struct {
	minimalAiexecVersion ti.Model

	warning string
}

func newVersionRequire() versionRequire {
	minimalAiexecVersion := ti.New()
	minimalAiexecVersion.Placeholder = "Minimal Aiexec version"
	minimalAiexecVersion.CharLimit = 128
	minimalAiexecVersion.Prompt = "Minimal Aiexec version (press Enter to next step): "
	minimalAiexecVersion.Focus()

	return versionRequire{
		minimalAiexecVersion: minimalAiexecVersion,
	}
}

func (p versionRequire) MinimalAiexecVersion() string {
	return p.minimalAiexecVersion.Value()
}

func (p versionRequire) View() string {
	s := fmt.Sprintf("Edit minimal Aiexec version requirement, leave it blank by default\n%s\n", p.minimalAiexecVersion.View())
	if p.warning != "" {
		s += fmt.Sprintf("\033[31m%s\033[0m\n", p.warning)
	}
	return s
}

func (p *versionRequire) checkRule() bool {
	if p.minimalAiexecVersion.Value() == "" {
		p.warning = ""
		return true
	}

	_, err := manifest_entities.NewVersion(p.minimalAiexecVersion.Value())
	if err != nil {
		p.warning = "Invalid minimal Aiexec version"
		return false
	}

	p.warning = ""
	return true
}

func (p versionRequire) Update(msg tea.Msg) (subMenu, subMenuEvent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return p, SUB_MENU_EVENT_NONE, tea.Quit
		case "enter":
			// check if empty
			if !p.checkRule() {
				return p, SUB_MENU_EVENT_NONE, nil
			}
			return p, SUB_MENU_EVENT_NEXT, nil
		}
	}

	// update view
	var cmd tea.Cmd
	p.minimalAiexecVersion, cmd = p.minimalAiexecVersion.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return p, SUB_MENU_EVENT_NONE, tea.Batch(cmds...)
}

func (p versionRequire) Init() tea.Cmd {
	return nil
}

func (p *versionRequire) SetMinimalAiexecVersion(version string) {
	p.minimalAiexecVersion.SetValue(version)
}
