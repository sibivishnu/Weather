@echo off
set NAME=webapp
set ONEBOX_IMAGE=onebox_webapp
set CONTAINER_PORT=5000:5000

if "%~1"=="build" goto build
if "%~1"=="compile" goto compile
goto default

:build
cp -r templates tmp
..\make.cmd build
goto done

:compile
cp -r templates tmp
..\make.cmd compile
goto done

:default
..\make.cmd %~1

:done
