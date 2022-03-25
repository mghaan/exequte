#Persistent
#NoEnv
#SingleInstance Ignore

procExe := 0

Menu, Tray, NoStandard
Menu, Tray, Icon, %A_ScriptDir%\exequte.ico
Menu, Tray, Add, exeQute, main
Menu, Tray, Disable, exeQute
Menu, Tray, Add
Menu, Tray, Add, Exit, handlerExit

main:
Run, %A_ScriptDir%\exequte.exe, %A_ScriptDir%, Hide, procExe
return

handlerExit:
Process, Close, %procExe%
ExitApp
return