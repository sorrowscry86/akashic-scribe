# Installation Guide

This document provides detailed instructions for installing Akashic Scribe on various platforms.

## System Requirements

- **Operating System**: Windows 10/11, macOS 10.14+, or Linux (X11/Wayland)
- **RAM**: 4GB minimum, 8GB recommended
- **Disk Space**: 500MB for application + processing space
- **Internet**: Required for URL-based video processing

## Windows Installation

1. Download the latest release from the [GitHub Releases page](https://github.com/sorrowscry86/akashic-scribe/releases)
2. Run the installer and follow the on-screen instructions
3. Launch Akashic Scribe from the Start Menu or desktop shortcut

## macOS Installation

1. Download the latest .dmg file from the [GitHub Releases page](https://github.com/sorrowscry86/akashic-scribe/releases)
2. Open the .dmg file and drag Akashic Scribe to your Applications folder
3. Launch from the Applications folder or Launchpad

## Linux Installation

### Debian/Ubuntu

```bash
# Add the VoidCat RDC repository
sudo add-apt-repository ppa:voidcat-rdc/akashic-scribe
sudo apt update

# Install Akashic Scribe
sudo apt install akashic-scribe
```

### Fedora/RHEL

```bash
# Add the VoidCat RDC repository
sudo dnf config-manager --add-repo https://repo.voidcat.org/akashic-scribe.repo
sudo dnf update

# Install Akashic Scribe
sudo dnf install akashic-scribe
```

### Arch Linux

```bash
# Using AUR helper (e.g., yay)
yay -S akashic-scribe
```

## Building from Source

If you prefer to build from source, follow these steps:

1. Ensure you have Go 1.24 or later installed
2. Clone the repository:
   ```bash
   git clone https://github.com/sorrowscry86/akashic-scribe.git
   ```
3. Navigate to the project directory:
   ```bash
   cd akashic-scribe/akashic_scribe
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Build the application:
   ```bash
   go build -o akashic_scribe .
   ```
6. Run the application:
   ```bash
   ./akashic_scribe
   ```

## Troubleshooting

If you encounter any issues during installation, please check our [Troubleshooting Guide](user-guide/troubleshooting.md) or [open an issue](https://github.com/sorrowscry86/akashic-scribe/issues) on GitHub.

## Next Steps

After installation, we recommend checking out the [User Manual](user-guide/index.md) to get started with Akashic Scribe.