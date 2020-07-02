@echo off
if "%OS%" === "Windows_NT" setolcal

echo wmic_get::model
wmic CPU get Name

echo wmic_get::logical_cpu
wmic CPU get NumberOfLogicalProcessors

echo wmic_get::core_cpu
wmic CPU get NumberOfCores

echo wmic_get::os
wmic OS get Name

echo wmic_get::MemTotal
wmic ComputerSystem get TotalPhysicalMemory

echo wmic_get::SwapTotal
wmic OS get TotalVirtualMemorySize
