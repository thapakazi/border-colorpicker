package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/lusingander/colorpicker"
)

// Config holds the user selections.
type Config struct {
	ActiveColor   string  `json:"active_color"`
	InactiveColor string  `json:"inactive_color"`
	BorderWidth   float64 `json:"border_width"`
}

const configFileName = "border_picker.json"

// configFilePath returns the path to the configuration file in ~/.config/.
func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", configFileName), nil
}

// loadConfig loads the configuration from the JSON file.
func loadConfig() (Config, error) {
	path, err := configFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// saveConfig writes the configuration to the JSON file.
func saveConfig(cfg Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	// Ensure the directory exists.
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// parseHexColor converts a hex string (in "0xAARRGGBB" format) to color.NRGBA.
func parseHexColor(s string) (color.NRGBA, error) {
	// Expect the string to have length 10 and start with "0x".
	if len(s) != 10 || s[:2] != "0x" {
		return color.NRGBA{}, fmt.Errorf("invalid color format: %s", s)
	}
	val, err := strconv.ParseUint(s[2:], 16, 32)
	if err != nil {
		return color.NRGBA{}, err
	}
	return color.NRGBA{
		A: uint8((val >> 24) & 0xFF),
		R: uint8((val >> 16) & 0xFF),
		G: uint8((val >> 8) & 0xFF),
		B: uint8(val & 0xFF),
	}, nil
}

func main() {
	// Load configuration from file (or use defaults if loading fails).
	cfg, err := loadConfig()
	if err != nil {
		log.Printf("Could not load config, using defaults: %v\n", err)
		cfg = Config{
			ActiveColor:   "0xffe2e2e3", // default active
			InactiveColor: "0xff414550", // default inactive
			BorderWidth:   6.0,
		}
	}
	activeColor, err := parseHexColor(cfg.ActiveColor)
	if err != nil {
		log.Printf("Error parsing active_color, using default: %v\n", err)
		activeColor = color.NRGBA{R: 0xe2, G: 0xe2, B: 0xe3, A: 0xff}
	}
	inactiveColor, err := parseHexColor(cfg.InactiveColor)
	if err != nil {
		log.Printf("Error parsing inactive_color, using default: %v\n", err)
		inactiveColor = color.NRGBA{R: 0x41, G: 0x45, B: 0x50, A: 0xff}
	}
	borderWidth := cfg.BorderWidth

	// Create the Fyne app and main window.
	a := app.New()
	w := a.NewWindow("Borders Color Picker")

	// Create displays for showing the current active and inactive colors.
	activeDisplay := newDisplay(activeColor)
	inactiveDisplay := newDisplay(inactiveColor)

	// Create the active color picker.
	activePicker := colorpicker.New(200, colorpicker.StyleHueCircle)
	activePicker.SetColor(activeColor)
	activePicker.SetOnChanged(func(c color.Color) {
		activeDisplay.setColor(c)
		activeColor = toNRGBA(c)
		updateBorders(activeColor, inactiveColor, borderWidth)
	})

	// Create the inactive color picker.
	inactivePicker := colorpicker.New(200, colorpicker.StyleHueCircle)
	inactivePicker.SetColor(inactiveColor)
	inactivePicker.SetOnChanged(func(c color.Color) {
		inactiveDisplay.setColor(c)
		inactiveColor = toNRGBA(c)
		updateBorders(activeColor, inactiveColor, borderWidth)
	})

	// Create a slider for border width.
	// The slider's range is from 0 to 20.
	widthSlider := widget.NewSlider(0, 20)
	widthSlider.Step = 0.1
	widthSlider.Value = borderWidth
	// A label to show the current width value.
	widthValueLabel := widget.NewLabel(fmt.Sprintf("%.1f", borderWidth))

	// Wrap the slider in a container and force it to be longer (e.g., 300 pixels wide).
	// Using container.NewWithoutLayout to prevent automatic layout resizing.
	sliderContainer := container.NewWithoutLayout(widthSlider)
	// Resize the slider to be 300 pixels wide.
	// (You might adjust the height as needed; here we keep its current minimal height.)
	widthSlider.Resize(fyne.NewSize(300, widthSlider.MinSize().Height))

	widthSlider.OnChanged = func(val float64) {
		borderWidth = val
		widthValueLabel.SetText(fmt.Sprintf("%.1f", borderWidth))
		updateBorders(activeColor, inactiveColor, borderWidth)
	}

	// --- Layout Setup ---

	// Active color section: title, then a horizontal container with the color picker and displays.
	activeSection := container.NewVBox(
		widget.NewLabel("Active Color"),
		container.NewHBox(
			activePicker,
			container.NewVBox(
				activeDisplay.label,
				activeDisplay.rect,
			),
		),
	)

	// Inactive color section.
	inactiveSection := container.NewVBox(
		widget.NewLabel("Inactive Color"),
		container.NewHBox(
			inactivePicker,
			container.NewVBox(
				inactiveDisplay.label,
				inactiveDisplay.rect,
			),
		),
	)

	// Border width section: title, then slider (in its container) and current value.
	widthSection := container.NewVBox(
		widget.NewLabel("Border Width"),
		container.NewHBox(
			sliderContainer,
			layout.NewSpacer(),
			widthValueLabel,
		),
	)

	// Compose the overall UI.
	content := container.NewVBox(
		activeSection,
		widget.NewSeparator(),
		inactiveSection,
		widget.NewSeparator(),
		widthSection,
	)
	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 800))
	w.ShowAndRun()
}

// updateBorders executes the borders binary with the new settings
// and then saves the current configuration.
func updateBorders(active, inactive color.NRGBA, width float64) {
	activeHex := hexColorString(active)
	inactiveHex := hexColorString(inactive)
	widthStr := fmt.Sprintf("width=%.1f", width)
	cmd := exec.Command("borders",
		"active_color="+activeHex,
		"inactive_color="+inactiveHex,
		widthStr)
	if err := cmd.Run(); err != nil {
		log.Printf("Error updating borders: %v\n", err)
	} else {
		fmt.Printf("Updated borders with active_color=%s, inactive_color=%s, %s\n",
			activeHex, inactiveHex, widthStr)
	}

	// Save current selections to config file.
	cfg := Config{
		ActiveColor:   activeHex,
		InactiveColor: inactiveHex,
		BorderWidth:   width,
	}
	if err := saveConfig(cfg); err != nil {
		log.Printf("Error saving config: %v\n", err)
	}
}

// display encapsulates a label and a rectangle to show a color sample.
type display struct {
	label *widget.Label
	rect  *canvas.Rectangle
}

// newDisplay creates a new display with the given initial color.
func newDisplay(clr color.Color) *display {
	nrgba := toNRGBA(clr)
	label := widget.NewLabel(hexColorString(nrgba))
	rect := canvas.NewRectangle(nrgba)
	rect.SetMinSize(fyne.NewSize(50, 50))
	return &display{
		label: label,
		rect:  rect,
	}
}

// setColor updates the display's label and rectangle to the given color.
func (d *display) setColor(clr color.Color) {
	nrgba := toNRGBA(clr)
	hex := hexColorString(nrgba)
	d.label.SetText(hex)
	d.rect.FillColor = nrgba
	d.rect.Refresh()
}

// hexColorString converts a color.NRGBA to a hex string in the format "0xAARRGGBB".
func hexColorString(c color.NRGBA) string {
	return fmt.Sprintf("0x%.2x%.2x%.2x%.2x", c.A, c.R, c.G, c.B)
}

// toNRGBA converts any color.Color to a color.NRGBA.
func toNRGBA(c color.Color) color.NRGBA {
	if nrgba, ok := c.(color.NRGBA); ok {
		return nrgba
	}
	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
