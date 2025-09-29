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
  - [Node.JS](#nodejs)
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

#### [Fish Shell](https://fishshell.com/)

<img src="https://fishshell.com/assets/img/Terminal_Logo_LCD_Small.png" width="17%" align="right" />

> Fish is a smart and user-friendly command lineshell for Linux, macOS, and the rest of the family.

Install `fish`:

```bash
brew instsall fish
```

Copy config:
```bash
mkdir -p $HOME/.config/fish/ && \
curl -L# -o $HOME/.config/fish/config.fish https://github.com/alexander-danilenko/dotfiles-macos/blob/main/files/.config/fish/config.fish
```

**[Oh My Fish](https://github.com/oh-my-fish/oh-my-fish)**: Package manager

> [Oh My Fish](https://github.com/oh-my-fish/oh-my-fish) provides core infrastructure to allow you to install packages which extend or modify the look of your shell. It's fast, extensible and easy to use..

Install `oh-my-fish`:

```bash
curl https://raw.githubusercontent.com/oh-my-fish/oh-my-fish/master/bin/install | fish
```

Run `fish` and install plugins:
```bash
omf install bobthefish bass nvm aws; omf theme bobthefish
```

## Development


<img src="https://cdn.svgporn.com/logos/nodejs.svg" width="17%" align="right" />

### Node.JS

> [`NVM`](https://github.com/nvm-sh/nvm) allows you to quickly install and use different versions of node via the command line.

Install NVM:

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash
```

Install LTS and set as default:

```bash
NODE_VERSION=20 && nvm install $NODE_VERSION && nvm alias default "$NODE_VERSION"
```

Install global packages:

```bash
NPM_PACKAGES=(
  '@nestjs/cli' # Nest.JS CLI
  contentful-cli # Contentful.com CLI
  dynamodb-admin # Handy Web-UI for viewing local DynamoDB data
  eslint
  eslint-plugin-import
  eslint-plugin-jsx-a11y
  eslint-plugin-node
  eslint-plugin-promise
  eslint-plugin-react
  eslint-plugin-react-hooks
  firebase-tools
  http-server # Simple HTTP server for static files in directory
  snyk # snyk.com CLI
  typescript
  @openapitools/openapi-generator-cli
) && npm install --global ${NPM_PACKAGES[@]}
```

<img src="https://cdn.svgporn.com/logos/visual-studio-code.svg" width="17%" align="right" />

### Visual Studio Code

Download config:
```bash
cd "$HOME/Library/Application Support/Code/User" && \
curl -O [$githubDir](https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/refs/heads/main/files/Library/Application%20Support/Code/User)/settings.json && \
curl -O [$githubDir](https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/refs/heads/main/files/Library/Application%20Support/Code/User)/keybindings.json && \
curl -O [$githubDir](https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/refs/heads/main/files/Library/Application%20Support/Code/User)/mcp.json
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
