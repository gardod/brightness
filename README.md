# brightness
Tray application for quickly changing monitor brightness

## Build

### Windows

To prevent launching a console window when running on Windows, add these command-line build flags:

```sh
go build -ldflags -H=windowsgui
```
