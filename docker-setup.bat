@echo off
REM ===========================================
REM Ezra API - Docker Quick Start Script
REM ===========================================
REM This script automates the Docker setup process
REM Run this after creating your .env file

setlocal enabledelayedexpansion

echo.
echo ============================================
echo   Ezra API - Docker Setup
echo ============================================
echo.

REM Check if running as administrator (optional but recommended)
net session >nul 2>&1
if %errorLevel% == 0 (
    echo [INFO] Running with administrator privileges
) else (
    echo [WARN] Not running as administrator
    echo        Some operations may require elevation
)

echo.
echo [STEP 1/8] Checking Docker...
echo --------------------------------------------

REM Check if Docker is installed
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not installed or not in PATH!
    echo.
    echo Please install Docker Desktop from:
    echo https://www.docker.com/products/docker-desktop/
    echo.
    pause
    exit /b 1
)

REM Check if Docker is running
docker ps >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker Desktop is not running!
    echo.
    echo Please start Docker Desktop and wait for it to fully start
    echo (Look for the whale icon in system tray)
    echo.
    echo Press any key after starting Docker Desktop...
    pause >nul
    
    REM Check again
    docker ps >nul 2>&1
    if errorlevel 1 (
        echo [ERROR] Docker still not responding
        echo Please ensure Docker Desktop is fully started
        pause
        exit /b 1
    )
)

echo [OK] Docker is installed and running
docker --version
docker-compose --version

echo.
echo [STEP 2/8] Checking .env file...
echo --------------------------------------------

if not exist .env (
    echo [WARN] .env file not found!
    echo.
    
    if exist .env.example (
        echo Creating .env from .env.example...
        copy .env.example .env >nul
        echo [OK] .env file created
        echo.
        echo ============================================
        echo   ACTION REQUIRED!
        echo ============================================
        echo.
        echo Please edit .env file and set at minimum:
        echo   1. DB_PASSWORD=your_secure_password
        echo   2. SECRET_KEY=your_secret_key_min_32chars
        echo.
        echo Opening .env in notepad...
        notepad .env
        echo.
        echo After saving .env, press any key to continue...
        pause >nul
    ) else (
        echo [ERROR] .env.example not found!
        echo Cannot create .env file automatically
        echo.
        echo Please create .env manually with required variables
        pause
        exit /b 1
    )
) else (
    echo [OK] .env file exists
)

REM Verify required environment variables
findstr /C:"DB_PASSWORD=" .env >nul
if errorlevel 1 (
    echo [WARN] DB_PASSWORD not set in .env
    set MISSING_VARS=1
)

findstr /C:"SECRET_KEY=" .env >nul
if errorlevel 1 (
    echo [WARN] SECRET_KEY not set in .env
    set MISSING_VARS=1
)

if defined MISSING_VARS (
    echo.
    echo [ERROR] Required environment variables are missing!
    echo Please edit .env and add:
    echo   - DB_PASSWORD
    echo   - SECRET_KEY
    echo.
    pause
    exit /b 1
)

echo.
echo [STEP 3/8] Checking network connectivity...
echo --------------------------------------------

REM Test Docker Hub connectivity
ping -n 1 registry-1.docker.io >nul 2>&1
if errorlevel 1 (
    echo [WARN] Cannot reach Docker Hub (registry-1.docker.io)
    echo.
    echo This might cause issues downloading images
    echo Possible solutions:
    echo   1. Check your internet connection
    echo   2. Check firewall settings
    echo   3. Try using a VPN
    echo   4. Use pre-downloaded images
    echo.
    echo Do you want to continue anyway? (Y/N)
    set /p CONTINUE=
    if /i not "!CONTINUE!"=="Y" (
        echo Aborting...
        exit /b 1
    )
) else (
    echo [OK] Network connectivity good
)

echo.
echo [STEP 4/8] Cleaning up old containers...
echo --------------------------------------------

docker-compose down >nul 2>&1
echo [OK] Cleanup complete

echo.
echo [STEP 5/8] Checking/pulling Docker images...
echo --------------------------------------------
echo This may take several minutes on first run...
echo.

REM Check if images exist, if not try to pull
docker image inspect postgres:16-alpine >nul 2>&1
if errorlevel 1 (
    echo [+] Pulling postgres:16-alpine...
    docker pull postgres:16-alpine
    if errorlevel 1 (
        echo [ERROR] Failed to pull postgres image
        echo.
        echo Possible solutions:
        echo   1. Check internet connection
        echo   2. Restart Docker Desktop
        echo   3. Try manual pull: docker pull postgres:16-alpine
        echo   4. Use a VPN if Docker Hub is blocked
        echo.
        pause
        exit /b 1
    )
) else (
    echo [OK] postgres:16-alpine image available
)

echo.
echo [STEP 6/8] Building application...
echo --------------------------------------------

docker-compose build
if errorlevel 1 (
    echo [ERROR] Build failed!
    echo Check the error messages above
    pause
    exit /b 1
)
echo [OK] Build complete

echo.
echo [STEP 7/8] Starting services...
echo --------------------------------------------

docker-compose up -d
if errorlevel 1 (
    echo [ERROR] Failed to start services!
    echo.
    echo Check docker-compose logs for details:
    echo   docker-compose logs
    echo.
    pause
    exit /b 1
)

echo [OK] Services started

echo.
echo [STEP 8/8] Waiting for services to be ready...
echo --------------------------------------------

REM Wait for database to be ready
set RETRY=0
:wait_db
set /a RETRY+=1
if %RETRY% GTR 30 (
    echo [WARN] Database took too long to start
    goto check_services
)

docker exec ezra-postgres pg_isready -U postgres >nul 2>&1
if errorlevel 1 (
    echo [%RETRY%/30] Waiting for database...
    timeout /t 2 >nul
    goto wait_db
)

echo [OK] Database is ready

:check_services
echo.
echo ============================================
echo   Setup Complete!
echo ============================================
echo.
echo Services Status:
echo --------------------------------------------
docker-compose ps

echo.
echo Service URLs:
echo --------------------------------------------
echo   PostgreSQL: localhost:5432
echo   API:        http://localhost:8080
echo.
echo Next Steps:
echo --------------------------------------------
echo   1. Wait 10-20 seconds for API to fully start
echo   2. Check API health: curl http://localhost:8080/health
echo   3. View logs: docker-compose logs -f
echo   4. Run migrations: make migrate-up
echo.
echo Useful Commands:
echo --------------------------------------------
echo   View logs:     docker-compose logs -f
echo   Stop services: docker-compose down
echo   Restart:       docker-compose restart
echo   Connect to DB: make db-connect
echo.

REM Try to check API health
timeout /t 5 >nul
echo Checking API health...
curl -f http://localhost:8080/health >nul 2>&1
if errorlevel 1 (
    echo [WARN] API not responding yet (normal on first start)
    echo Wait 30 seconds and try: curl http://localhost:8080/health
) else (
    echo [OK] API is healthy!
)

echo.
echo Press any key to view logs (Ctrl+C to exit logs)...
pause >nul

docker-compose logs -f

endlocal
