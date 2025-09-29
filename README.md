<p align="center"><img src="https://cdn.svgporn.com/logos/macOS.svg" /></p>

<h1 align="center">MacOS Dotfiles and system config</h1>

📖 Table of contents:
- [Prerequisites](#prerequisites)
- [Software](#software)
  - [Quick Setup with Brewfile](#quick-setup-with-brewfile)
  - [Terminal](#terminal)
    - [Git](#git)
    - [ZSH](#zsh)
    - [Fish Shell](#fish-shell)
- [Development](#development)
  - [Visual Studio Code](#visual-studio-code)
- [Containerization](#containerization)
  - [Docker](#docker)
  - [Docksal](#docksal)
- [Hardware](#hardware)
  - [Apple keyboard](#apple-keyboard)
    - [Karabiner](#karabiner)
- [Troubleshooting](#troubleshooting)


---

## Prerequisites

Brew: https://brew.sh/
  ```shell
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  ```


This will install all the essential CLI tools, GUI applications, development tools, and fonts in one command. The Brewfile includes:

## Software

### Quick Setup with Brewfile

For a streamlined installation of all packages and applications, use the provided [Brewfile](./files/Brewfile):

```bash
# Download the Brewfile
curl -O https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/Brewfile

# Install all packages and applications
brew bundle
```


### Terminal

#### Git

```shell
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.gitconfig > ~/.gitconfig && \
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.gitignore > ~/.gitignore
```

#### ZSH

```shell
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.profile >> ~/.profile
```

#### [ZSH with Oh My Zsh](https://ohmyz.sh/)

<img src="https://ohmyz.sh/img/OMZLogo_BnW.png" width="17%" align="right" />

> Oh My Zsh is a delightful, open source, community-driven framework for managing your Zsh configuration.

Install Oh My Zsh:
```bash
sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

Copy ZSH config:
```bash
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.zshrc > ~/.zshrc
```

**Oh My Zsh Features**:
- Over 200+ plugins and themes
- Auto-completion and syntax highlighting
- Git integration with useful aliases
- Plugin ecosystem for development tools
- Easy customization and theming

**Included Plugins**:
- `git` - Git aliases and functions
- `aws` - AWS CLI completion
- `docker` - Docker aliases and completion
- `nvm` - Node Version Manager integration
- `python` - Python development tools
- `brew` - Homebrew integration
- `macos` - macOS-specific utilities
- `vscode` - VS Code integration
- `zsh-autosuggestions` - Command suggestions
- `zsh-syntax-highlighting` - Syntax highlighting

## Development

<img src="https://cdn.svgporn.com/logos/visual-studio-code.svg" width="17%" align="right" />

### Visual Studio Code

Download config:
```bash
cd "$HOME/Library/Application Support/Code/User" && \
curl -O https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/refs/heads/main/files/Library/Application%20Support/Code/User/settings.json -O https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/refs/heads/main/files/Library/Application%20Support/Code/User/keybindings.json -O https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/refs/heads/main/files/Library/Application%20Support/Code/User/mcp.json
```

> **Note**: VS Code extensions are now managed via the [Brewfile](#quick-setup-with-brewfile) and will be installed automatically with `brew bundle`.

## Containerization

<img src="https://cdn.svgporn.com/logos/docker.svg" width="17%" align="right" />

### Docker

https://docs.docker.com/desktop/install/mac-install/

<img src="https://avatars.githubusercontent.com/u/20954974" width="17%" align="right" />

### Docksal

https://docksal.io/installation#linux-supported

```shell
bash <(curl -fsSL https://get.docksal.io)
```

## Hardware

### Apple keyboard

#### Karabiner

Install [Karabiner](https://karabiner-elements.pqrs.org/):
```shell
brew install --cask karabiner-elements
```

- [Fix Ukrainian Apple Keyboard rule](./files/.config/karabiner/rules/fix-ukrainian-apple-keyboard.json)


## Troubleshooting

TBD
