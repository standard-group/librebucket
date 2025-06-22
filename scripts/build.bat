@echo off
echo Building Librebucket for Windows...

echo Building Go binary...
 go build -o librebucket.exe
if %errorlevel% neq 0 (
    echo Error: Go build failed
    exit /b %errorlevel%
)

echo Build completed successfully!