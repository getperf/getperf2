管理用単体移行調査
GPFConfig, GPFSetupConfig, GPFSetupOption 調査
管理用 SOAP 部 UML 作成
管理用 SOAP 部 TODO 作成

# getperf2 調査

Module        
ElapseTime    
StartTime     
Mode          
ManagedPid    
LocaleFlag    
DaemonFlag    
Host          
ServiceName   
Pwd           
Home          
ParameterFile 
ProgramName   
ProgramPath   

OutDir        
WorkDir       
WorkCommonDir 
ArchiveDir    
BackupDir     
ScriptDir     
BinDir        
LogDir        

SslDir       
CacertFile   
ClcertFile   
ClkeyFile    
SvParamFile  
SvcertFile   
SvkeyFile    
SvcacertFile 
LicenseFile  

SoapRetry

ExitFlag 
PidFile  
PidPath  

SslConfig
Schedule

# agent

ReserveSender( filename )
SendData( filename)
SendMessage( severity, message )
DownloadCertificate( long timestamp )
ReserveFileSender( onOff, waitSec )  
SendZipData( filename)
DownloadConfigFilePM( filename )

# admin

GetLatestBuild(  )
RegistAgent( GPFSetupConfig *setup )
DownloadUpdateModule( int _build, char *moduleFile)  
CheckHostStatus( GPFSetupConfig *setup)



