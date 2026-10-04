<#
.SYNOPSIS
  Run all AoC day solutions sequentially for perf testing.

.DESCRIPTION
  Iterates over day* folders (natural order: day2 before day10), builds and
  runs each one. Each day's own perf timer prints as usual. This script adds
  outer headers, a summary table (with each day's parsed Average time), and a
  grand average across all days.

  benchmark.md only ever keeps the best (lowest) time per day. If a new run
  is slower than or equal to the stored value, the old value is kept.

.EXAMPLE
  .\run-all.ps1                         # run all except the default excludes
  .\run-all.ps1 -Exclude day10,day12    # exclude several
  .\run-all.ps1 -Exclude @()            # run everything
  .\run-all.ps1 -NoBuild                # use `go run .` instead of building
  pwsh ./2015/run-all.ps1 -Exclude day10

.NOTES
  Requires Go installed and `go.mod` at repo root.
  Expects each day's program to print a line like:
    Average: 12.3456 microseconds
  (ms / ns / us units are also understood and converted to microseconds.)
  This file is ASCII-only on purpose, so it parses correctly regardless of encoding.
#>
param(
    [string]$Year = "2015",
    [string[]]$Exclude = @("day60"),
    [string]$BenchmarkFile = $null,   # defaults to "$BaseDir\benchmark.md"
    [switch]$NoBuild,                 # fall back to `go run .`
    [switch]$DisableBoundsCheck,
    [switch]$KeepExe
)

# === Edit this to permanently exclude folders (used when -Exclude isn't passed) ===
$DefaultExclude = @("day60")
if (-not $PSBoundParameters.ContainsKey('Exclude')) {
    $Exclude = $DefaultExclude
}

$inv = [System.Globalization.CultureInfo]::InvariantCulture
$isWin = ($PSVersionTable.PSVersion.Major -lt 6) -or $IsWindows

# Resolve base dir (target year folder)
$BaseDir = $PSScriptRoot
if (-not $BaseDir) { $BaseDir = (Get-Location).Path }
if (Test-Path (Join-Path $BaseDir $Year)) {
    $BaseDir = Join-Path $BaseDir $Year
}

$ExcludeLower = @($Exclude | ForEach-Object { $_.ToLowerInvariant().Trim() } | Where-Object { $_ -ne "" })

if (-not $BenchmarkFile) {
    $BenchmarkFile = Join-Path $BaseDir "benchmark.md"
}

function Get-DayNum([string]$n) {
    $m = [regex]::Match($n, '\d+')
    if ($m.Success) { [int]$m.Value } else { 0 }
}

function ConvertTo-Micro([double]$v, [string]$unit) {
    $u = $unit.ToLowerInvariant()
    if ($u -eq 'ms' -or $u -like 'milli*') { return $v * 1000 }
    if ($u -eq 'ns' -or $u -like 'nano*')  { return $v / 1000 }
    return $v
}

Write-Host "Base dir : $BaseDir" -ForegroundColor DarkGray
if ($ExcludeLower.Count -gt 0) {
    Write-Host "Excluding: $($ExcludeLower -join ', ')" -ForegroundColor Yellow
} else {
    Write-Host "Excluding: (none)" -ForegroundColor Green
}
Write-Host ""

$dayDirs = @(Get-ChildItem -Path $BaseDir -Directory -Filter "day*" |
    Sort-Object { Get-DayNum $_.Name }, Name)

if ($dayDirs.Count -eq 0) {
    Write-Error "No day* folders found in $BaseDir"
    exit 1
}

# "Average: 12.3456 microseconds" (ANSI colour codes are stripped before matching)
$AvgRegex  = [regex]'(?im)^\s*Average:\s*([\d.]+)\s*(microseconds?|[\u00B5\u03BC]s|us|milliseconds?|ms|nanoseconds?|ns)\b'
$AnsiRegex = [regex]'\x1B\[[0-9;?]*[A-Za-z]'

$results  = @()
$globalSw = [System.Diagnostics.Stopwatch]::StartNew()

foreach ($dir in $dayDirs) {
    $name      = $dir.Name
    $nameLower = $name.ToLowerInvariant()

    if ($ExcludeLower -contains $nameLower) {
        Write-Host "=== $name [SKIPPED - excluded] ===" -ForegroundColor DarkGray
        $results += [PSCustomObject]@{ Day = $name; Status = "SKIPPED"; Elapsed = "-"; AvgUs = $null }
        continue
    }

    if (-not (Get-ChildItem -Path $dir.FullName -Filter "*.go" -ErrorAction SilentlyContinue)) {
        Write-Host "=== $name [SKIPPED - no .go file] ===" -ForegroundColor DarkGray
        $results += [PSCustomObject]@{ Day = $name; Status = "NO GO FILE"; Elapsed = "-"; AvgUs = $null }
        continue
    }

    Write-Host "=== $name ===" -ForegroundColor Cyan

    $exeName = if ($isWin) { "$name.exe" } else { $name }
    $exePath = Join-Path $dir.FullName $exeName
    $sw      = [System.Diagnostics.Stopwatch]::new()
    $avgUs   = $null

    Push-Location $dir.FullName
    try {
        if (-not $NoBuild) {
            $buildArgs = @("build", "-ldflags=-s -w", "-trimpath")
            if ($DisableBoundsCheck) { $buildArgs += '-gcflags=all=-B' }
            $buildArgs += @("-o", $exeName, ".")
            $buildOutput = & go @buildArgs 2>&1
            if ($LASTEXITCODE -ne 0) {
                throw "Build failed: $(($buildOutput | Out-String).Trim())"
            }
            $sw.Start()
            $output = & $exePath 2>&1
            $exitCode = $LASTEXITCODE
            $sw.Stop()
        } else {
            $sw.Start()
            $output = & go run . 2>&1
            $exitCode = $LASTEXITCODE
            $sw.Stop()
        }

        # Normalise (stderr lines arrive as ErrorRecords in Windows PowerShell)
        $lines = @($output | ForEach-Object { "$_" })
        $lines | ForEach-Object { Write-Host $_ }

        # Parse the day's self-reported average, normalised to microseconds
        $clean = $AnsiRegex.Replace(($lines -join "`n"), '')
        $m = $AvgRegex.Match($clean)
        if ($m.Success) {
            $avgUs = ConvertTo-Micro ([double]::Parse($m.Groups[1].Value, $inv)) $m.Groups[2].Value
        }

        if ($exitCode -eq 0) {
            Write-Host "[OK] $name finished in $($sw.ElapsedMilliseconds)ms (outer)" -ForegroundColor Green
            $status = "OK"
        } else {
            Write-Host "[FAIL] $name failed (exit $exitCode) after $($sw.ElapsedMilliseconds)ms" -ForegroundColor Red
            $status = "FAIL ($exitCode)"
        }
        $results += [PSCustomObject]@{
            Day = $name; Status = $status; Elapsed = "$($sw.ElapsedMilliseconds)ms"; AvgUs = $avgUs
        }
    } catch {
        $sw.Stop()
        Write-Host "[FAIL] $name exception: $_" -ForegroundColor Red
        $results += [PSCustomObject]@{
            Day = $name; Status = "EXCEPTION"; Elapsed = "$($sw.ElapsedMilliseconds)ms"; AvgUs = $null
        }
    } finally {
        Pop-Location
        if (-not $KeepExe -and -not $NoBuild -and (Test-Path $exePath)) {
            Remove-Item $exePath -Force -ErrorAction SilentlyContinue
        }
    }
    Write-Host ""
}

$globalSw.Stop()

$display = $results | ForEach-Object {
    [PSCustomObject]@{
        Day     = $_.Day
        Status  = $_.Status
        Elapsed = $_.Elapsed
        Avg     = if ($null -ne $_.AvgUs) { "$($_.AvgUs.ToString('0.0000', $inv)) us" } else { "-" }
    }
}

Write-Host "=======================================" -ForegroundColor Cyan
Write-Host " Summary ($($results.Count) days, total outer $($globalSw.Elapsed.TotalSeconds.ToString('0.00', $inv))s)" -ForegroundColor Cyan
$display | Format-Table -AutoSize | Out-String | Write-Host

# Only successful runs count towards averages and the benchmark file
$withAvg = @($results | Where-Object { $_.Status -eq "OK" -and $null -ne $_.AvgUs })

if ($withAvg.Count -gt 0) {
    $grandAvg = ($withAvg | Measure-Object -Property AvgUs -Average).Average
    $fastest  = $withAvg | Sort-Object AvgUs | Select-Object -First 1
    $slowest  = $withAvg | Sort-Object AvgUs -Descending | Select-Object -First 1

    Write-Host "---------------------------------------" -ForegroundColor Cyan
    Write-Host " Average across $($withAvg.Count) days: $($grandAvg.ToString('0.0000', $inv)) microseconds" -ForegroundColor Magenta
    Write-Host " Fastest: $($fastest.Day) ($($fastest.AvgUs.ToString('0.0000', $inv)) us)" -ForegroundColor Green
    Write-Host " Slowest: $($slowest.Day) ($($slowest.AvgUs.ToString('0.0000', $inv)) us)" -ForegroundColor Yellow
} else {
    Write-Host "No days reported an 'Average: X microseconds' line to aggregate." -ForegroundColor DarkYellow
}

# --- Write / update benchmark.md (keeps the best time per day) ---
if ($withAvg.Count -gt 0) {
    $best = @{}
    if (Test-Path $BenchmarkFile) {
        foreach ($line in (Get-Content -Path $BenchmarkFile)) {
            if ($line -match '^\|\s*(day\w+)\s*\|\s*([\d,]+(?:\.\d+)?)\s*\|\s*$') {
                $best[$Matches[1]] = [double]::Parse($Matches[2].Replace(',', ''), $inv)
            }
        }
    }

    foreach ($r in $withAvg) {
        if ($best.ContainsKey($r.Day) -and $best[$r.Day] -le $r.AvgUs) {
            Write-Host "$($r.Day): kept previous best $($best[$r.Day].ToString('N4', $inv)) us (this run: $($r.AvgUs.ToString('N4', $inv)) us)" -ForegroundColor DarkGray
        } else {
            $best[$r.Day] = $r.AvgUs
        }
    }

    $mdLines  = @("# Benchmark Results", "", "| Day | Average Time ($([char]0x03BC)s) |", "|---|---:|")
    foreach ($day in ($best.Keys | Sort-Object { Get-DayNum $_ }, { $_ })) {
        $mdLines += "| $day | $($best[$day].ToString('N4', $inv)) |"
    }

    $sumUs    = ($best.Values | Measure-Object -Sum).Sum
    $mdLines += "| **Total** | **$($sumUs.ToString('N4', $inv)) ($(($sumUs / 1000).ToString('0.00', $inv)) ms)** |"

    $mdLines | Set-Content -Path $BenchmarkFile -Encoding utf8
    Write-Host ""
    Write-Host "Benchmark table written to: $BenchmarkFile" -ForegroundColor Cyan
} else {
    Write-Host ""
    Write-Host "Skipped writing benchmark.md - no averages to report." -ForegroundColor DarkYellow
}

Write-Host ""
Write-Host "Edit `$DefaultExclude at the top of run-all.ps1 to change the default excludes." -ForegroundColor DarkGray