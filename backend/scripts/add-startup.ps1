param(
    [string]$AppPath = "C:\temp\portarius.exe",
    [string]$AppName = "Portarius"
)

$WshShell = New-Object -ComObject WScript.Shell
$StartupFolder = [System.Environment]::GetFolderPath('Startup')
$ShortcutPath = Join-Path $StartupFolder "$AppName.lnk"
$Shortcut = $WshShell.CreateShortcut($ShortcutPath)
$Shortcut.TargetPath = $AppPath
$Shortcut.WorkingDirectory = Split-Path $AppPath
$Shortcut.Save()

Write-Host "Atalho criado em $ShortcutPath"
Write-Host "O aplicativo $AppName foi adicionado à inicialização automática do Windows."