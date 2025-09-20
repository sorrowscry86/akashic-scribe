# Installation Guide

**VoidCat RDC - Akashic Scribe**

This guide provides detailed instructions for installing and setting up Akashic Scribe on different operating systems.

## Table of Contents

- [System Requirements](#system-requirements)
- [Pre-built Binaries](#pre-built-binaries)
- [Building from Source](#building-from-source)
- [Troubleshooting](#troubleshooting)
- [Platform-Specific Notes](#platform-specific-notes)
- [Support](#support)

## System Requirements

### Minimum Requirements

- **Operating System**: Windows 10/11, macOS 10.14+, or Linux (with X11 or Wayland)
- **RAM**: 4GB minimum
- **Disk Space**: 500MB for the application
- **Internet Connection**: Required for URL-based video processing and updates
- **Graphics**: Hardware-accelerated graphics support recommended

### Recommended Requirements

- **RAM**: 8GB or more
- **CPU**: Multi-core processor for optimal performance
- **Disk Space**: 2GB+ for processing large video files
- **Graphics**: Dedicated GPU for enhanced performance

## Pre-built Binaries

### Windows

#### Option 1: Windows Installer (Recommended)

1. Download the latest Windows installer from the [Releases](https://github.com/sorrowscry86/akashic-scribe/releases) page
2. Run the installer with administrator privileges
3. Follow the installation wizard
4. Launch Akashic Scribe from the Start menu or desktop shortcut

#### Option 2: Portable Version

1. Download the portable ZIP file from the [Releases](https://github.com/sorrowscry86/akashic-scribe/releases) page
2. Extract to your preferred directory
3. Run `akashic_scribe.exe` directly

### macOS

#### Option 1: DMG Installer (Recommended)

1. Download the latest macOS .dmg file from the [Releases](https://github.com/sorrowscry86/akashic-scribe/releases) page
2. Open the .dmg file
3. Drag Akashic Scribe to your Applications folder
4. Launch from Applications folder or Launchpad
5. If prompted about unknown developer, go to System Preferences > Security & Privacy and allow the app

#### Option 2: Homebrew

```bash
# Add VoidCat RDC tap
brew tap voidcat/akashic-scribe

# Install Akashic Scribe
brew install akashic-scribe
```

### Linux

#### Debian/Ubuntu

```bash
# Add the VoidCat RDC repository
curl -fsSL https://repo.voidcat.com/gpg | sudo apt-key add -
echo "deb [arch=amd64] https://repo.voidcat.com/akashic-scribe stable main" | sudo tee /etc/apt/sources.list.d/akashic-scribe.list

# Update package list
sudo apt update

# Install Akashic Scribe
sudo apt install akashic-scribe
```

#### Fedora/RHEL/CentOS

```bash
# Add the VoidCat RDC repository
sudo dnf config-manager --add-repo https://repo.voidcat.com/akashic-scribe.repo

# Install Akashic Scribe
sudo dnf install akashic-scribe
```

#### Arch Linux (AUR)

```bash
# Using yay
yay -S akashic-scribe-bin

# Or using paru
paru -S akashic-scribe-bin
```

#### AppImage (Universal)

```bash
# Download AppImage
wget https://github.com/sorrowscry86/akashic-scribe/releases/latest/download/akashic-scribe-x86_64.AppImage

# Make executable
chmod +x akashic-scribe-x86_64.AppImage

# Run
./akashic-scribe-x86_64.AppImage
```

## Building from Source

### Prerequisites

- **Go**: Version 1.24.4 or later
- **Git**: For cloning the repository
- **GCC**: Compatible C compiler
- **Development Tools**: Platform-specific build tools

### Platform-Specific Prerequisites

#### Windows

```powershell
# Install Go
# Download from https://golang.org/dl/

# Install Git
# Download from https://git-scm.com/

# Install TDM-GCC or MinGW-w64
# Download from https://jmeubank.github.io/tdm-gcc/
```

#### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Go and dependencies via Homebrew
brew install go gcc
```

#### Linux

##### Debian/Ubuntu

```bash
# Install development tools
sudo apt update
sudo apt install golang gcc git build-essential libgl1-mesa-dev xorg-dev

# For Wayland support
sudo apt install libwayland-dev libxkbcommon-dev
```

##### Fedora/RHEL/CentOS

```bash
# Install development tools
sudo dnf install golang gcc git libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel

# For Wayland support
sudo dnf install wayland-devel libxkbcommon-devel
```

##### Arch Linux

```bash
# Install development tools
sudo pacman -S go gcc git libgl libxcursor libxrandr libxi libxinerama

# For Wayland support
sudo pacman -S wayland libxkbcommon
```

### Build Process

```bash
# Clone the repository
git clone https://github.com/sorrowscry86/akashic-scribe.git
cd akashic-scribe/akashic_scribe

# Install Go dependencies
go mod tidy

# Build for current platform
go build -o akashic_scribe .

# Cross-compilation examples
# For Windows from other platforms
GOOS=windows GOARCH=amd64 go build -o akashic_scribe.exe .

# For macOS from other platforms
GOOS=darwin GOARCH=amd64 go build -o akashic_scribe .

# For Linux from other platforms
GOOS=linux GOARCH=amd64 go build -o akashic_scribe .
```

### Installation

```bash
# Optional: Install to system PATH
# Linux/macOS
sudo cp akashic_scribe /usr/local/bin/

# Windows (as Administrator)
copy akashic_scribe.exe "C:\Program Files\VoidCat\AkashicScribe\"
```

## Troubleshooting

### Common Issues

#### Application Fails to Start

**Symptoms**: Application crashes on startup or fails to launch

**Solutions**:

- Ensure you have the latest version of Go installed (if building from source)
- Check that all dependencies are installed correctly
- Verify your system meets the minimum requirements
- Run from terminal/command prompt to see error messages
- Check antivirus software isn't blocking the application

#### Missing GUI Elements

**Symptoms**: Application window appears but controls are missing or malformed

**Solutions**:

- Update your graphics drivers to the latest version
- Ensure you have the required X11 or Wayland libraries installed (Linux)
- Try running with different display settings
- Check if hardware acceleration is supported

#### Permission Issues

**Symptoms**: Unable to save files or access certain directories

**Solutions**:

- Run the application with appropriate permissions
- Ensure write access to the output directory
- Check file system permissions
- Try running as administrator (Windows) or with sudo (Linux/macOS) temporarily

#### Network Connectivity Issues

**Symptoms**: URL-based video processing fails

**Solutions**:

- Check your internet connection
- Verify firewall settings allow the application
- Test with a different video URL
- Check proxy settings if behind corporate firewall

#### Performance Issues

**Symptoms**: Application runs slowly or uses excessive resources

**Solutions**:

- Ensure you meet the recommended system requirements
- Close other resource-intensive applications
- Check available disk space
- Monitor CPU and memory usage
- Consider processing smaller video files

### Platform-Specific Issues

#### Windows

- **Issue**: "MSVCP140.dll not found"
  - **Solution**: Install Microsoft Visual C++ Redistributable 2015-2022

- **Issue**: Windows Defender flags the application
  - **Solution**: Add application to Windows Defender exclusions

#### macOS

- **Issue**: "App can't be opened because it is from an unidentified developer"
  - **Solution**: Right-click app, select Open, then click Open in dialog

- **Issue**: Application appears in dock but window doesn't show
  - **Solution**: Check Mission Control or try Alt+Tab to find window

#### Linux

- **Issue**: "cannot open display" error
  - **Solution**: Ensure DISPLAY environment variable is set correctly

- **Issue**: Wayland compatibility issues
  - **Solution**: Try running with X11 backend or install additional Wayland libraries

### Getting Help

If you encounter issues not covered in this guide:

1. **Search Existing Issues**: Check the [Issue Tracker](https://github.com/sorrowscry86/akashic-scribe/issues)
2. **Community Support**: Visit our [Discussions](https://github.com/sorrowscry86/akashic-scribe/discussions)
3. **Create New Issue**: Provide detailed information including:
   - Operating system and version
   - Application version
   - Steps to reproduce the issue
   - Error messages or logs
   - System specifications

## Platform-Specific Notes

### Windows Notes

- Windows 10 version 1903 or later recommended for best compatibility
- Windows Defender SmartScreen may require approval on first run
- UWP version available in Microsoft Store (coming soon)

### macOS Notes

- macOS Big Sur (11.0) or later recommended for optimal performance
- Apple Silicon (M1/M2) native builds available
- Notarized builds ensure security compliance

### Linux Notes

- Tested on Ubuntu 20.04+, Fedora 35+, Arch Linux
- Both X11 and Wayland display servers supported
- Flatpak version available for universal compatibility

## Updating

### Automatic Updates

Akashic Scribe includes an automatic update checker that:

- Checks for updates on startup
- Notifies you when updates are available
- Can download and install updates (with permission)

### Manual Updates

To manually update Akashic Scribe:

1. Download the latest version from the [Releases](https://github.com/sorrowscry86/akashic-scribe/releases) page
2. Install over your existing installation
3. Your settings and preferences will be preserved

### Development Builds

For the latest features and fixes:

```bash
# Clone and build latest development version
git clone https://github.com/sorrowscry86/akashic-scribe.git
cd akashic-scribe/akashic_scribe
git checkout develop
go mod tidy
go build -o akashic_scribe .
```

---

## 📞 Support & Contact

- **GitHub Issues**: [Report bugs or request features](https://github.com/sorrowscry86/akashic-scribe/issues)
- **Discussions**: [Community discussions and Q&A](https://github.com/sorrowscry86/akashic-scribe/discussions)
- **Developer**: [@sorrowscry86](https://github.com/sorrowscry86)
- **Project**: VoidCat RDC - Akashic Scribe
- **Contact**: Wykeve Freeman (Sorrow Eternal) - SorrowsCry86@voidcat.org
- **Organization**: VoidCat RDC
- **Support Development**: CashApp $WykeveTF

---

**© 2024 VoidCat RDC, LLC. All rights reserved.**