# macOS Dotfiles | AI Agent Rules

## PROJECT OVERVIEW
macOS dotfiles repository with Homebrew Bundle automation, VS Code configuration, and shell setup.

## BREWFILE MANAGEMENT
- **Install**: `brew bundle` (installs all packages from Brewfile)
- **Structure**: Organized sections (CLI tools, GUI apps, VS Code extensions, system utilities)
- **Organization**: Maintain existing categorization, use descriptive comments, group related packages

### AI Instructions for Brewfile
When working with this project's Brewfile, AI agents should:
1. **Understand the structure**: The Brewfile uses Homebrew Bundle syntax with organized sections for CLI tools, GUI applications, development tools, and system utilities
2. **Respect the organization**: Maintain the existing categorization and commenting structure when adding new packages
3. **Follow naming conventions**: Use descriptive comments and group related packages together
4. **Consider dependencies**: Some packages require additional setup steps (noted in comments)
5. **Use proper syntax**: Follow Homebrew Bundle syntax for different package types

### Brewfile Syntax Reference
```ruby
# Homebrew formulae (CLI tools)
brew "package-name"

# Homebrew casks (GUI applications)  
cask "application-name"

# VS Code extensions
vscode "publisher.extension-name"

# Mac App Store apps
mas "App Name", id: 123456789

# Homebrew taps (third-party repositories)
tap "user/repo"
```

### Mac App Store Dependencies
When adding Mac App Store apps to the Brewfile:
1. **Extract app information from App Store URLs**:
   - URL format: `https://apps.apple.com/ua/app/{app-name}/id{app-id}?mt=12`
   - Extract app name from URL path (e.g., `amphetamine` from `/app/amphetamine/`)
   - Extract numeric app ID from URL (e.g., `937984704` from `/id937984704`)
2. **Add to Brewfile using `mas` syntax**: `mas "App Name", id: 123456789`
3. **Example transformation**:
   - Input URL: `https://apps.apple.com/ua/app/amphetamine/id937984704?mt=12`
   - Brewfile entry: `mas "amphetamine", id: 937984704`
4. **Placement**: Add Mac App Store apps in the appropriate section (usually GUI Applications)

## VS CODE CONFIGURATION
- **Extensions**: Managed via Brewfile (vscode entries)
- **Settings**: JetBrains Mono font, GitHub Dark theme, material icons
- **Keybindings**: Custom navigation (cmd+[/] for back/forward)
- **MCP**: Atlassian integration configured

## GIT CONFIGURATION
- **Aliases**: `c` (commit), `s` (status), `p` (pull with rebase), `history` (graph log)
- **Tools**: Meld for diff/merge, nano editor
- **Settings**: Auto CRLF input, long paths enabled

## PATTERNS
- **Brewfile**: Group by category with clear section headers
- **VS Code**: Extensions installed via Homebrew Bundle
- **Shell**: Fish with Oh My Fish, ZSH with custom profile
- **Git**: Rebase-based workflow with custom aliases

## SECURITY
- Standard dotfiles patterns (no sensitive data)
- Git configuration includes safe defaults
- VS Code settings use standard security practices

## GIT CONVENTIONS
- **Commit Format**: Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification
- **Structure**: `<type>[optional scope]: <description>`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`, `build`
- **Breaking Changes**: Use `!` after type/scope or `BREAKING CHANGE:` footer
- **Examples**: 
  - `feat: add new VS Code extension`
  - `fix(brewfile): correct package name`
  - `docs: update installation instructions`
  - `feat!: remove deprecated configuration`
- **Workflow**: Rebase-based (`p` alias pulls with rebase)
- **Aliases**: Custom aliases for common operations