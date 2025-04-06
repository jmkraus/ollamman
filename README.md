# OllaMMan - The ollama Model Manager

## Description

The ollama Model Manager (OllaMMan) is a commandline tool to extend the "list" command of the excellent local LLM runner [ollama](https://ollama.com).

Instead of `ollama list` you can now call `ollamman` for a list of locally installed LLM's. It basically provides the same information, but has some additional options:

| Parameter           | Description                                  |
| ------------------- | -------------------------------------------- |
| -d, --order-date    | orders models by last modified date          |
| -n, --order-name    | orders by name (ascending)                   |
| -c, --check-updates | checks for each model if update is available |

### Update feedback

 * Checkmark (green): latest version - no updates available
 * Exclamation (yellow): a newer version of the model is available for download
 * X mark (red): web site not found (e.g. if it is a custom model)
 * Circle (gray): no update check

*Disclaimer:* This app has been hacked together quite quickly and didn't receive much testing, yet. It serves my purposes quite well, and may be of some use to others. But don't please blame me if it fries your computer!

## Installation

To compile, you should have the GoLang package installed on your machine.

1. Checkout the repository
2. Install the dependencies
```
go mod init github.com/jmkraus/ollamman
go mod tidy  
```
3. Compile & run

It has been tested with v0.6.4 of ollama on macOS (Sequoia) only and builds successfully with Go 1.24.1.

## License

Distributed under MIT License

JMK 2025
