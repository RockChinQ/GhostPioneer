cd D:\download\go1.15.windows-386\go\bin
go.exe build -i -ldflags -H=windowsgui
copy main.exe ..\gl-32bit.exe
pause