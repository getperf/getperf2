# getperfsoapメイン移行

変更調査
プロトタイプ

# 変更調査

既存ソース調査

getperfsoap.c

```c
char *gpfHelpMessage[] = {
    "getperfsoap [--send(-s)|--get(-g)] [--config(-c) getperf.cfg]",
    "            filename.zip",
    "Options:",
    "  -s --send                 send data",
    "  -g --get                  get data",
    "  -c --config <getperf.cfg> config file",
    "  <filename.zip>",
    0 /* end of text */
};

```

プロキシの設定チェック、

```c
    /* If PROXY_ENABLE is true and PROXY_HOST is if blank (NULL), Proxy settings to apply environment variable; HTTP_PROXY. */
    /* It PROXY_ENABLE is false, Proxy setting is disable even if an environment variable has been set. */
    gpfCheckHttpProxyEnv( &(config->schedule) );

```

メイン処理

```c
        if ( sendFlag == 1 )
        {
            for ( retry = 0; retry < GPF_SOAP_RETRY; retry ++ )
            {
                rc = gpfReserveSender( config, zipfile );
                if ( rc == 1 )
                {
                    if ( ( rc = gpfSendData( config, zipfile ) ) == 1 )
                        gpfNotice("[Sended] %s", zipfile);
                    break;
                }

                gpfError("send data failed retry %d/%d", retry +1, GPF_SOAP_RETRY);
                if ( waitSec > 0 )
                {
                    sleep( waitSec );
                }
            }
        }
        else if ( getFlag == 1 )
        {
            for ( retry = 0; retry < GPF_SOAP_RETRY; retry ++ )
            {
                if ( ( rc = gpfDownloadCertificate( config, timestamp ) ) == 1 )
                {
                    gpfNotice("[Saved] %s", zipfile);
                    break;
                }
                gpfError("get data failed retry %d/%d", retry + 1, GPF_SOAP_RETRY);
                sleep( waitSec );
            }
        }

```


# TODO

実行引数パーサー[]
送信処理[]
受信処理[]
プロキシー設定[]

; --------- HTTP proxy, timeout setting --------------------------------
; The presence or absence of proxy settings
PROXY_ENABLE = false

; Proxy server address (in the case of the blank is applied HTTP_PROXY environment variable)
PROXY_HOST =

; Proxy server port
PROXY_PORT = 8080

; HTTP connection time-out
SOAP_TIMEOUT = 300

