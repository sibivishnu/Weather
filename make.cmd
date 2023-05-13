@echo off
setlocal enabledelayedexpansion

set IMAGE=azopat/gomagik
set "ROOT=%cd%"
set VENDOR_ROOT=%ROOT%\..\_vendor
set COMMON_ROOT=%ROOT%\..\common
set IMAGE_BASE_NAME=%DOCKER_IMAGES_BASE_NAME%/%NAME%
set LOCAL_ROOT=/go/%NAME%

set "DOCKER_OPTIONS=-v %ROOT%:%LOCAL_ROOT% -v %VENDOR_ROOT%:/go/_vendor -v %COMMON_ROOT%:/go/common"
set DOCKER_CMD_TO_RUN=docker run --rm=true %DOCKER_OPTIONS% %IMAGE%

rem set defaults
if "%COMMIT_HASH%"=="" (
    git show -s --format=%%h > commit_hash.tmp
    set /p COMMIT_HASH=<commit_hash.tmp
    rm commit_hash.tmp
)
if "%COMMIT_TIME%"=="" (
    git show -s --format=%%at > commit_time.tmp
    set /p COMMIT_TIME=<commit_time.tmp
    rm commit_time.tmp
)
if "%REDIS_HOST%"=="" set REDIS_HOST=127.0.0.1:6379
if "%ACCU_API_KEY%"=="" set ACCU_API_KEY="SET_ACCU_API_KEY_ENV_VAR"

if "%~1"=="run" goto run
if "%~1"=="attach" goto attach
if "%~1"=="build" goto build
if "%~1"=="compile" goto compile
if "%~1"=="gofmt" goto gofmt
if "%~1"=="details" goto details
if "%~1"=="help" goto help
goto help

:run
  echo "Running %ONEBOX_IMAGE%"
  docker run --rm=true -e ACCU_API_KEY=%ACCU_API_KEY% -e REDIS_HOST=%REDIS_HOST% -p %CONTAINER_PORT% --name %ONEBOX_IMAGE% %ONEBOX_IMAGE%
goto done

:attach
docker run -t -i --rm=true %DOCKER_OPTIONS% -p %CONTAINER_PORT% -e ACCU_API_KEY=%ACCU_API_KEY% -e REDIS_HOST=%REDIS_HOST% %IMAGE% /bin/bash
goto done


:build
echo "Building GO executable for %NAME%"
cp Win.Dockerfile tmp/Dockerfile
cp %ROOT%\..\auth\credential_file.json tmp
rm -rf tmp/%NAME%
%DOCKER_CMD_TO_RUN%  sh -c "cd /go && rm -rf /go/%NAME%/bin && GOPATH=/go/_vendor /usr/local/go/bin/go build -ldflags '-X main.BUILD=%COMMIT_TIME%__%COMMIT_HASH%'  -o  /go/%NAME%/tmp/%NAME% /go/%NAME%/*.go"
echo "Building Docker Image: %ONEBOX_IMAGE%"
docker build tmp -t %ONEBOX_IMAGE%
goto done

:compile
gofmt build
cp Win.Dockerfile tmp/Dockerfile
cp %ROOT%\..\auth\credential_file.json tmp
goto done

:gofmt
"%DOCKER_CMD_TO_RUN% sh -c "cd /go && rm -rf /go/%NAME%/bin && /usr/local/go/bin/gofmt -w /go/%NAME%/*.go"
goto done

:details
echo NAME=%NAME%
echo CONTAINER_PORT=%CONTAINER_PORT%
echo REDIS_HOST=%REDIS_HOST
echo COMMIT_TIME=%COMMIT_TIME%
echo COMMIT_HASH=%COMMIT_HASH%
echo LOCAL_ROOT=%LOCAL_ROOT%
echo VENDOR_ROOT=%VENDOR_ROOT%
echo COMMON_ROOT=%COMMON_ROOT%
echo ONEBOX_IMAGE=%ONEBOX_IMAGE%
echo DOCKER=%DOCKER_CMD_TO_RUN%
goto done

:help
echo "make.cmd attach|build|compile|gofmt (From dir of project)"
goto done

:done
