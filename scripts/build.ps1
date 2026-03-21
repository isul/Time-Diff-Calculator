#requires -Version 5.1
<#
.SYNOPSIS
  timediff 프로덕션 빌드 (기본: 테스트 후 모든 주요 OS용 wails build)
.DESCRIPTION
  기본 -platform:
    - macOS 호스트: windows/amd64,linux/amd64,darwin/amd64,darwin/arm64
    - Windows/Linux 호스트: windows/amd64,linux/amd64
  -Native: 현재 Windows만 빌드 (wails 기본 동작)
  -AllPlatforms: macOS 타깃 포함 강제 시도
  -WailsArgs에 -platform이 있으면 그대로 사용 (단일/복수 지정)
.EXAMPLE
  .\scripts\build.ps1
  .\scripts\build.ps1 -Native
  .\scripts\build.ps1 -NoTest
  .\scripts\build.ps1 -Clean
  .\scripts\build.ps1 -AllPlatforms
  .\scripts\build.ps1 -WailsArgs @('-platform','windows/amd64')
#>
param(
    [switch]$NoTest,
    [switch]$Clean,
    [switch]$Native,
    [switch]$AllPlatforms,
    [string[]]$WailsArgs = @()
)

$ErrorActionPreference = 'Stop'
$isMacOSHost = $PSVersionTable.OS -match 'Darwin|macOS'
if ($isMacOSHost) {
    $DefaultPlatforms = 'windows/amd64,linux/amd64,darwin/amd64,darwin/arm64'
} else {
    $DefaultPlatforms = 'windows/amd64,linux/amd64'
}

$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

function Test-HasPlatformArg {
    param([string[]]$Args)
    for ($i = 0; $i -lt $Args.Count; $i++) {
        if ($Args[$i] -eq '-platform' -or $Args[$i] -like '-platform=*') {
            return $true
        }
    }
    return $false
}

if ($Clean) {
    $binDir = Join-Path $Root 'build\bin'
    if (Test-Path $binDir) {
        Remove-Item -Recurse -Force $binDir
    }
    Write-Host '[build] cleaned build/bin'
}

if (-not $NoTest) {
    Write-Host '[build] go test ./...'
    go test ./...
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}

$hasPlat = Test-HasPlatformArg $WailsArgs

if ($AllPlatforms) {
    $DefaultPlatforms = 'windows/amd64,linux/amd64,darwin/amd64,darwin/arm64'
    Write-Host '[build] -AllPlatforms: macOS 타깃 포함 시도'
    if (-not $isMacOSHost) {
        Write-Host '[build] 경고: 현재 OS에서는 macOS 타깃이 Wails 제약으로 실패/건너뛰기 될 수 있습니다.'
    }
}

if ($Native -or $hasPlat) {
    Write-Host "[build] wails build $($WailsArgs -join ' ')"
    if ($WailsArgs.Count -gt 0) {
        & wails build @WailsArgs
    } else {
        & wails build
    }
} else {
    Write-Host "[build] wails build -platform $DefaultPlatforms $($WailsArgs -join ' ')"
    & wails build -platform $DefaultPlatforms @WailsArgs
}

if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "[build] 완료: $(Join-Path $Root 'build\bin')"
