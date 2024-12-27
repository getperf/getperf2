既存コード調査
UML 更新

# 既存コード調査

## getperfsoap.c

main()関数内

```c
if ( sendFlag == 1 )

    gpfCheckHttpProxyEnv( &(config->schedule) );
    rc = gpfReserveSender( config, zipfile );
    if ( ( rc = gpfSendData( config, zipfile ) ) == 1 )
        break;

else if ( getFlag == 1 )

    if ( ( rc = gpfDownloadCertificate( config, timestamp ) ) == 1 )
```

## gpf_admin.c

```c
gpfSetUserInfo( GPFSetupOption *options )
gpfUpdateModule( GPFConfig *config, int build, char *archive, int forkFlag )
gpfRunCheckCoreUpdate( GPFConfig *config )
gpfEntryHost( GPFConfig *config, GPFSetupConfig *setup )
gpfRunSetup( GPFSetupOption *options )
gpfDeployConfigFile( GPFConfig *config, char *pass, char *configFile)
```

gpfSetUserInfo( options )
gpfUpdateModule( config, build, archive, fork )   
gpfRunCheckCoreUpdate( config )
gpfEntryHost( config, setup )
gpfRunSetup( options )
gpfDeployConfigFile( config, pass, configFile)

## gpf_agent.c

```c
gpfExecSOAPCommandPM( GPFConfig *config, char *option, char *filePath )

int gpfCheckLicense( GPFConfig *config, int expiredTime )

    if (! gpfExecSOAPCommandPM( config, "--get" , "sslconf.zip" ) )
    else if (! gpfUnzipSSLConf( config ) )

int gpfSendCollectorDataAll( GPFTask *task )

    if ( ( rc = gpfSendCollectorData( task, target ) ) == 0 )

int gpfSendCollectorData( GPFTask *task, char *zipFile )

    rc = gpfExecSOAPCommandPM( config, "--send", zipFile ) ;

```

gpfExecSOAPCommandPM( config, option, filePath )  
gpfCheckLicense( config, expiredTime )
gpfSendCollectorDataAll( task )
gpfSendCollectorData( task, zipFile )

getperfsoap.cを移行するのが良さそう。

## gpf_soap_admin.c

```c
gpsetGetLatestBuild( GPFConfig *config )
gpsetRegistAgent( GPFConfig *config, GPFSetupConfig *setup )
gpfDownloadUpdateModule( GPFConfig *config, int _build, char *moduleFile)
gpfCheckHostStatus( GPFConfig *config, GPFSetupConfig *setup)

```

## gpf_soap_agent.c

```c
gpfReserveSender( GPFConfig *config, char *filename )
gpfSendData( GPFConfig *config, char *filename)
gpfSendMessage( GPFConfig *config, int severity, char *message )
gpfDownloadCertificate( GPFConfig *config, long timestamp )
gpfReserveFileSender( GPFConfig *config, char *onOff, int *waitSec )
gpfSendZipData( GPFConfig *config, char *filename)
gpfDownloadConfigFilePM( GPFConfig *config, char *filename )

```


## gpf_soap_common.c

```c
_gpfSoapSSLInit()
gpfSetSoapProperties(struct soap *soap, GPFSchedule *schedule)
gpfSoapSSLClientContext( GPFConfig *config, struct soap *soap, int sslType)
soap *gpfCreateSoap(GPFConfig *config, int sslType)
soap *gpfCreateSoapWithMime(
    GPFConfig *config, int sslType, char *zipBuffer, size_t zipSize, char *filename )
gpfSoapError(struct soap *soap, char *msg)
gpfReadZipFile( GPFConfig *config, const char *zipPath, size_t *_zipSize)
gpfGetFileFromMIME( GPFConfig *config, struct soap *soap, char *filename)

```


