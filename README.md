# :construction: [NURU](https://nuruprogramming.org/) LSP SERVER :construction:

## Still under development :smile:

## For more info on Nuru Programming Language click [here](https://nuruprogramming.org/)

## Current Features:

- Variables and Function completions
- Go to definitions
- Hover on functions, keywords and variables
- Standard Library completions
- Cross pakeji completions
- Error highlighting (Doesn't work in latest VScode, supported in almost every other editor)

## Building from source

```sh
git clone https://github.com/Borwe/nuru-lsp
cd nuru-lsp
go mod tidy
go build
```

## Setting up LSP with your editor

Once you build you will have the executable `nuru-lsp`, you can use it to setup with your editor and relate it with `.nr` or `.sr` files.

#### Example VsCode

- It is adviced to [this README.md](https://github.com/Borwe/nuru-lsp/blob/master/vscode-plugin/README.md) for more information
- Also for better UX install nuru-snippets by Kabwoy and the official Nuru Programming Vscode plugin

#### Example for neovim:

```lua
-- requires you to have nvim-lspconfig
local lspconfig = require("lsp-config")
local lspconfig_configurer = require('lspconfig.configs')
lspconfig_configurer["nuru-lsp"] = {
  default_config = {
    cmd = { '/Path/to/nuru-lsp' },
    filetypes = { 'sr', 'nroff' },
    root_dir = require('lspconfig.util').find_git_ancestor,
    single_file_support = true,
  },
  docs = {
    description = [[
https://github.com/Borwe/nuru-lsp

Nuru Unofficial Language Server
        ]],
    default_config = {
      root_dir = lspconfig.util.find_git_ancestor]],
    },
  },
}
```

#### Example for helix:

Inside `languages.toml`

```toml
[language-server.nuru-lsp]
command = "/path/to/nuru-lsp"

[[language]]
name = "Nuru"
scope = "source.nr"
file-types = ["nr", "sr"]
comment-token = "//"
indent = { tab-width = 4, unit = " "}
language-servers = ["nuru-lsp"]
```

#### Example Emacs with Eglot:

```elisp
;; Nuru-LSP
(define-derived-mode nuru-mode prog-mode "Nuru Mode")
(setq auto-mode-alist
	  (append '(("\\nr\\'" . nuru-mode)
				("\\sr\\'" . nuru-mode))
			  auto-mode-alist))
(add-hook 'nuru-mode-hook 'eglot-ensure)
;;Add path to nuru on emacs search path
(add-to-list 'exec-path "/Path/To/nuru-lsp")
(add-to-list 'eglot-server-programs
			 '(nuru-mode . ("nuru-lsp")))
```
