windres -o main.syso icon.rc

go build -ldflags="-s -w -H=windowsgui" -o="W:\Software\CloudMusic\网易云音乐.exe"