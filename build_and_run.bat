@echo off
echo Setting up Fyne development environment...

REM Add MinGW-w64 and Go to PATH
set PATH=C:\msys64\ucrt64\bin;C:\Program Files\Go\bin;%PATH%

REM Enable CGO for Fyne
set CGO_ENABLED=1
set CC=gcc

echo Environment configured.
echo Building the application...
cd akashic_scribe
go build -v .

echo Running the application...
.\akashic_scribe.exe
