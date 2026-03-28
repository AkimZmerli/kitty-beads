# VHS Integration Guide

## What is VHS?

[VHS](https://github.com/charmbracelet/vhs) is a terminal GIF recorder and player from Charmbracelet. It allows you to create terminal recordings and convert them to GIFs, MP4s, or WebM videos programmatically using a simple script syntax.

## Installation

```bash
# macOS
brew install charmbracelet/tap/vhs

# Linux
docker pull ghcr.io/charmbracelet/vhs
# or build from source
go install github.com/charmbracelet/vhs@latest
```

## Integration into Kitty Beads

### Setup
1. Add VHS config directory: `docs/vhs/` for `.tape` scripts
2. Configure output directory: `docs/recordings/` for generated GIFs/videos
3. Use `make record` or similar target to generate recordings

### Basic Tape Script Example
```tape
# Create a simple demo recording
Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 600

Type "bd list --tui"
Enter
Sleep 2
Key ctrl+c
```

## Use Cases for Kitty Beads

### 1. **Feature Demonstrations**
   - Create GIFs showing TUI in action (list view, graph, status)
   - Record daemon mode operations and workflows
   - Demonstrate filter syntax and keyboard navigation

### 2. **Documentation & Onboarding**
   - Animated guides in README for new users
   - Feature overview GIFs in docs/
   - Step-by-step workflow recordings

### 3. **Issue & PR Documentation**
   - Attach recordings to issues showing bugs
   - Include before/after recordings in PRs
   - Create visual changelogs for releases

### 4. **Testing & CI/CD**
   - Automated regression testing via frame comparison
   - Generate documentation during CI pipeline
   - Compare recordings across versions

### 5. **Command Reference**
   - Record all major commands: `bd list --tui`, `bd graph --tui`, `bd status --tui`
   - Create interactive help guides
   - Show daemon interactions and monitoring

## Integration with Bubble Tea TUI

VHS pairs well with the existing Bubble Tea TUI implementation:
- Record actual `bd list --tui` sessions
- Capture Tokyo Night theme styling
- Show real keyboard interactions and animations

## Example Workflow

```bash
# Create recordings
vhs < docs/vhs/list-demo.tape

# Output files
docs/recordings/list-demo.gif
docs/recordings/list-demo.mp4

# Embed in docs
![bd list --tui demo](../recordings/list-demo.gif)
```

## Future Enhancements
- Automate recording generation in CI for all TUI commands
- Create interactive demo in web interface
- Generate changelog GIFs for releases
