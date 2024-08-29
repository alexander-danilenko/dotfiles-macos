if status is-interactive
  set fish_greeting "" # Disables greeting message.
end

###> `bobthefish` Theme configuration
function fish_greeting; end
set -g theme_color_scheme dracula
set -g theme_display_date no
set -g theme_display_node yes
set -g theme_newline_cursor yes
set -g theme_newline_prompt '$ '
set -g theme_title_use_abbreviated_path no
###< `bobthefish` Theme configuration