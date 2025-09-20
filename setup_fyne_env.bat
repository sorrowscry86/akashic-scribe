@echo off
echo Setting up Fyne development environment...

REM Add MinGW-w64 and Go to PATH
set PATH=C:\msys64\ucrt64\bin;C:\Program Files\Go\bin;%PATH%

REM Enable CGO for Fyne
set CGO_ENABLED=1
set CC=gcc

echo Environment configured:
echo - MinGW-w64 GCC: C:\msys64\ucrt64\bin
echo - Go: C:\Program Files\Go\bin  
echo - CGO_ENABLED: %CGO_ENABLED%
echo - CC: %CC%
echo.

echo Testing tools...
gcc --version
echo.
go version
echo.
go env CGO_ENABLED

echo.
echo Ready for Fyne development!
echo.
echo To build Akashic Scribe:
echo   cd akashic_scribe
echo   go build -v .
echo.
echo To run:
echo   .\akashic_scribe.exe
