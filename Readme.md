# JankyBorder in Steroids with colorpicker

<img align="right" width="50%" src="images/picker.png" alt="Borders Color Picker">

This is a Fyne-based GUI tool written in Go that lets you dynamically update the border properties of the [JankyBorders](https://github.com/FelixKratz/JankyBorders) binary.

With this tool you can adjust:

  - **Active Color:** The color of the active border.
  - **Inactive Color:** The color of the inactive border.
  - **Border Width:** The thickness of the border (with a range from 0 to 20).

    Changes are immediately applied to the running `borders` instance, and your settings are saved to a configuration file so they persist between sessions.


## Features

- **Dynamic Updates:** Immediately apply changes by launching the `borders` binary with new options.
- **Dual Color Pickers:** Separate color selectors for active and inactive border colors.
- **Border Width Slider:** Adjust the border width with an extended slider (0 to 20, with the slider resized for ease-of-use).
- **Configuration Persistence:** Saves your current selections to `~/.config/border_picker.json` and loads them on startup.

<img align="right" width="50%" src="images/demo.png" alt="Borders Color Picker Demo">


## Prerequisites

- [Go](https://golang.org/) (latest greatest, whatever, mine: go1.23.5)
- [Fyne](https://fyne.io/) (GUI framework for Go)
- [lusingander/colorpicker](https://github.com/lusingander/colorpicker)
- The `borders` binary from [JankyBorders](https://github.com/FelixKratz/JankyBorders) must be available in your PATH.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/thapakazi/border-colorpicker.git
   cd border-colorpicker
   go get
   go run main.go
   ```

## Usage
First ensure the borders binary is available in your PATH.
```
go build  -o bcolorpik main.go
./bcolorpik
```

Your changes are applied in real time by re-launching the borders binary with updated options, and your selections are saved to ~/.config/border_picker.json for future sessions.

### Configuration File

The configuration file is stored at ~/.config/border_picker.json and holds your current selections in JSON format. An example configuration file looks like:
```js
{
  "active_color": "0xffe2e2e3",
  "inactive_color": "0xff414550",
  "border_width": 6.0
}
```

## Thanks 🙇🏼 
- [Fyne](https://fyne.io/)
- [lusingander/colorpicker](https://github.com/lusingander/colorpicker)
- [JankyBorders](https://github.com/FelixKratz/JankyBorders)

