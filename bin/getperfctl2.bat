@echo off
echo %0
for %%F in (%0) do set cwd=%%~dpF

set HOME=%cwd%..

%HOME%\bin\getperf2.exe  %1 %2 %3 %4 %5 %6 %7 %8 %9

