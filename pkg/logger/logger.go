package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)
)

func Success(message string) {
	fmt.Println(successStyle.Render("✅ " + message))
}

func Error(message string) {
	fmt.Println(errorStyle.Render("❌ " + message))
}

func Info(message string) {
	fmt.Println(infoStyle.Render("📋 " + message))
}

func Warning(message string) {
	fmt.Println(warningStyle.Render("⚠️  " + message))
}

func Step(icon, message string) {
	fmt.Println(infoStyle.Render(icon + " " + message))
}

func Println(message string) {
	fmt.Println(message)
}
