# soap_agent_admin 単体

RegistAgent
DownloadUpdateModule
CheckHostStatus

# RgistAgent

リグレッション

go test ./agent/ --run TestRegistAgent  -v

# DownloadUpdateModule

リグレッション

go test ./agent/ --run TestDownloadUpdateModule -v

# CheckHostStatus

リグレッション

go test ./agent/ --run TestCheckHostStatus -v

