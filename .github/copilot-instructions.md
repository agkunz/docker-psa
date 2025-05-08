<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Docker PSA Project

This project is a Docker CLI plugin written in Go that enhances the output of `docker ps` with a more human-readable format.

Key files:
- `cmd/main.go`: Contains the main implementation of the plugin
- `plugin.json`: The Docker CLI plugin configuration
- `Makefile`: Contains build and installation commands

When suggesting code or improvements, consider:
1. Maintaining human readability of the output format
2. Following Docker plugin best practices
3. Optimizing performance when listing containers

When testing:
1. Use the `Makefile` to build the plugin, and then if it succeeds run `docker psa` right away