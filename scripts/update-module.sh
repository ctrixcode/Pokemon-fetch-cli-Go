#!/usr/bin/env bash
set -euo pipefail

# Derive module path from git remote 'origin'
origin_url=$(git remote get-url origin 2>/dev/null || true)
if [[ -z "$origin_url" ]]; then
  echo "No git remote 'origin' found. Aborting." >&2
  exit 1
fi

# Normalize to https form
# Supported examples:
#  git@github.com:owner/repo.git
#  https://github.com/owner/repo.git
#  https://github.com/owner/repo

url_no_scheme=${origin_url#git@}
url_no_scheme=${url_no_scheme#https://}
url_no_scheme=${url_no_scheme#ssh://git@}
url_no_scheme=${url_no_scheme#http://}
url_no_scheme=${url_no_scheme%/.git}
url_no_scheme=${url_no_scheme%.git}

# Convert git@github.com:owner/repo to github.com/owner/repo
url_no_scheme=${url_no_scheme/:/\/}

if [[ $url_no_scheme != github.com/* ]]; then
  echo "Remote is not a GitHub repository (got: $origin_url). Aborting." >&2
  exit 1
fi

module_path=$url_no_scheme

current_module_line=$(grep -E '^module ' go.mod || true)
if [[ -z "$current_module_line" ]]; then
  echo "go.mod missing module line. Aborting." >&2
  exit 1
fi

current_module=${current_module_line#module }

if [[ "$current_module" == "$module_path" ]]; then
  echo "Module already up to date: $module_path"
  exit 0
fi

echo "Updating module path: $current_module -> $module_path"
# In-place replace module line (portable enough for GNU/BSD sed variants)
# Try BSD sed first fallback to GNU
if sed -i '' -e "1s|.*|module $module_path|" go.mod 2>/dev/null; then
  :
else
  sed -i -e "1s|.*|module $module_path|" go.mod
fi

go mod tidy

echo "Module path updated and dependencies tidied."
