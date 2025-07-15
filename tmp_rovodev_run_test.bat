@echo off
echo Running direct test...
"C:\Program Files\Go\bin\go.exe" run tmp_rovodev_direct_test.go
echo.
echo Exit code: %ERRORLEVEL%
echo.
echo Press any key to continue...
pause > nul