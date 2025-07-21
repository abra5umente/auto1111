@echo off
SETLOCAL

REM Change to project root
cd /d %~dp0

REM Optional: activate your virtual environment
REM call venv\Scripts\activate.bat

REM Run the FastAPI server
echo Starting backend server...
uvicorn backend.main:app --reload

ENDLOCAL
pause
