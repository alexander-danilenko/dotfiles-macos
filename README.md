<h1 align="center">MacOS Dotfiles and system config</h1>

## Prerequisites

Brew: https://brew.sh/
  ```shell
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  ```

## Software

Musthave non-gui packages:

```bash
BREW_CLI_PACKAGES=(
  jq
  yq
  mc
  ffmpeg
  yt-dlp
  fish
) && brew install ${BREW_CLI_PACKAGES[@]} 
```

Musthave gui apps:

```bash
BREW_GUI_PACKAGES=(
  1password
  anydesk
  zoom
  slack
  font-jetbrains-mono
  jetbrains-toolbox
  synology-drive
  google-chrome
  firefox@developer-edition
  postman
) && brew install --quiet --casks ${BREW_GUI_PACKAGES[@]} 
```

### Terminal

#### Git

```shell
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.gitconfig > ~/.gitconfig && \
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.gitignore > ~/.gitignore
```

#### ZSH

```shell
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.zprofile >> ~/.zprofile
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

### Node.JS

> [`NVM`](https://github.com/nvm-sh/nvm) allows you to quickly install and use different versions of node via the command line.

Install NVM:

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash
```

Install LTS and set as default:

```bash
NODE_VERSION=20 && \
nvm install $NODE_VERSION && \
nvm alias default "$NODE_VERSION"
```

Install global packages:

```bash
NPM_PACKAGES=(
  '@nestjs/cli' # Nest.JS CLI
  contentful-cli # Contentful.com CLI
  dynamodb-admin # Handy Web-UI for viewing local DynamoDB data
  eslint
  eslint-config-airbnb
  eslint-config-google
  eslint-config-standard
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
) && npm install --global ${NPM_PACKAGES[@]}
```

<img src="https://cdn.svgporn.com/logos/php.svg" width="17%" align="right" />

<img src="https://cdn.svgporn.com/logos/visual-studio-code.svg" width="17%" align="right" />

### Visual Studio Code

https://code.visualstudio.com/docs/setup/mac


```bash
brew install --cask visual-studio-code
```

Download config:
```bash
curl https://raw.githubusercontent.com/alexander-danilenko/dotfiles-macos/main/files/.config/Code/User/settings.json > ~/.config/Code/User/settings.json
```

Install extensions

```bash
CODE_EXTENSIONS=(
  GitHub.copilot                       # GitHub Copilot!
  GitHub.vscode-github-actions         # Github Actions support
  acarreiro.calculate                  # Calculates inline math expr
  christian-kohler.path-intellisense   # File path autocomplete
  dakara.transformer                   # Filter, Sort, Unique, Reverse, Align, CSV, Line Selection, Text Transformations and Macros
  dotenv.dotenv-vscode                 # .env support
  editorconfig.editorconfig            # EditorConfig support
  golang.go                            # Golang support
  ms-azuretools.vscode-docker          # Docker support
  ms-python.python                     # Python support
  ms-vscode-remote.remote-ssh          # SSH support
  tommasov.hosts                       # Hosts file syntax highlighter
  tyriar.lorem-ipsum                   # Lorem Ipsum generator
  yzhang.markdown-all-in-one           # Markdown tools
  #redhat.ansible                      # Ansible support

  # Node/NPM/Yarn specific extensions
  christian-kohler.npm-intellisense # NPM better autocomplete
  dbaeumer.vscode-eslint            # Eslint support
  
  # Themes
  github.github-vscode-theme    # GitHub color theme
  pkief.material-icon-theme     # Material Icon Theme
  rokoroku.vscode-theme-darcula # JetBrains-like theme
) && for extension in "${CODE_EXTENSIONS[@]}"; do
  code --install-extension "$extension" --force
done
```

## Cloud

<img src="https://cdn.svgporn.com/logos/aws.svg" width="17%" align="right" />

### AWS Cli

https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

<img src="https://cdn.svgporn.com/logos/google-cloud.svg" width="17%" align="right" />

### Google Cloud

https://cloud.google.com/sdk/docs/install#mac

<img src="https://raw.githubusercontent.com/opentofu/brand-artifacts/main/full/transparent/SVG/on-light.svg" width="17%" align="right" />

### OpenTofu 

```shell
brew install opentofu
```

<img src="https://cdn.svgporn.com/logos/terraform-icon.svg" width="17%" align="right" />

### Terraform 

```shell
brew install terraform
```

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


## Troubleshooting

TBD