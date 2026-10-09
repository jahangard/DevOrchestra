!include "MUI2.nsh"
Name "DevOrchestra"
OutFile "dist\DevOrchestra-Setup.exe"
InstallDir "$LOCALAPPDATA\Programs\DevOrchestra"
RequestExecutionLevel user
ShowInstDetails show

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Section "DevOrchestra (required)" SEC_APP
 SectionIn RO
 SetOutPath "$INSTDIR"
 File /oname=DevOrchestra.exe "dist\devorchestra-windows-amd64.exe"
 CreateDirectory "$SMPROGRAMS\DevOrchestra"
 CreateShortcut "$SMPROGRAMS\DevOrchestra\DevOrchestra.lnk" "$INSTDIR\DevOrchestra.exe"
 CreateShortcut "$DESKTOP\DevOrchestra.lnk" "$INSTDIR\DevOrchestra.exe"
 WriteUninstaller "$INSTDIR\Uninstall.exe"
 WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DevOrchestra" "DisplayName" "DevOrchestra"
 WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DevOrchestra" "UninstallString" "$INSTDIR\Uninstall.exe"
SectionEnd

Section /o "Install Codex CLI (requires existing Node.js and npm)" SEC_CODEX
 DetailPrint "Installing @openai/codex through npm..."
 nsExec::ExecToLog 'cmd.exe /C "npm.cmd install -g @openai/codex"'
 Pop $0
 StrCmp $0 "0" installed
 MessageBox MB_ICONEXCLAMATION "Could not install Codex CLI automatically. Install Node.js first, then run: npm install -g @openai/codex. DevOrchestra can also retry this from Settings."
 installed:
SectionEnd

Section "Uninstall"
 Delete "$DESKTOP\DevOrchestra.lnk"
 Delete "$SMPROGRAMS\DevOrchestra\DevOrchestra.lnk"
 RMDir "$SMPROGRAMS\DevOrchestra"
 Delete "$INSTDIR\DevOrchestra.exe"
 Delete "$INSTDIR\Uninstall.exe"
 RMDir "$INSTDIR"
 DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DevOrchestra"
 ; User data and Codex CLI are intentionally NOT removed.
SectionEnd
