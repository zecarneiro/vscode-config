pushd %1
cd /d %~dp0
@echo off
powershell.exe .\make.ps1 %*
popd