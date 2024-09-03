####### Node Version Manager ###################################################
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion

######## Aliases ###############################################################
alias node-modules-ls='find . -name "node_modules" -type d -prune -print | xargs du -chs'
alias node-modules-rm='find . -name "node_modules" -type d -prune -exec rm -rf "{}" +'
alias ll="ls -la"

######## PATH variable overriding ##############################################
PATH_DIRS=(
  "$HOME/bin"
  "$HOME/.local/bin"
  "$HOME/.composer/vendor/bin"
  "$HOME/Projects/bin"
  "$HOME/miniconda3/bin"
  "$HOME/Projects/bin/google-cloud-sdk/bin"
  "/opt/homebrew/opt/openjdk/bin"
)
for PATH_DIR in "${PATH_DIRS[@]}"; do
  if [ -d "$PATH_DIR" ] && [[ ":$PATH:" != *":$PATH_DIR:"* ]]; then
    PATH="$PATH_DIR:$PATH"
  fi
done

