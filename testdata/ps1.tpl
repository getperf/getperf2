Param(
    [string]$log_dir
)
$log_dir = Convert-Path $log_dir

{{range $i, $v := .}}
echo TestId::{{$v.TestId}}
$log_path = Join-Path $log_dir "{{$v.TestId}}"
Invoke-Command  -ScriptBlock { 
    {{$v.Script}} 
} | Out-File $log_path -Encoding UTF8{{end}}
