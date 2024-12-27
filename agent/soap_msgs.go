package agent

// テンプレート
var SoapRequestMsgTemplates map[string]string = map[string]string{
	"getLatestBuild": `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:ns0="http://perf.getperf.com">
    <soapenv:Header>
        <wsa:Action>urn:getLatestBuild</wsa:Action>
        <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
        <wsa:To>{{ .ToURL }}</wsa:To>
    </soapenv:Header>
    <soapenv:Body>
        <ns0:getLatestBuild xmlns:ns0="http://perf.getperf.com" xmlns="http://rmi.java/xsd">
            <ns0:moduleTag>{{ .moduleTag }}</ns0:moduleTag>
            <ns0:majorVer>{{ .majorVer }}</ns0:majorVer>
        </ns0:getLatestBuild xmlns:ns0="http://perf.getperf.com">
    </soapenv:Body>
</soapenv:Envelope>
`,
  "downloadCertificate": `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:ns0="http://perf.getperf.com">
    <soapenv:Header>
        <wsa:Action>urn:downloadCertificate</wsa:Action>
        <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
        <wsa:To>{{ .ToURL }}</wsa:To>
    </soapenv:Header>
    <soapenv:Body>
        <ns0:downloadCertificate xmlns:ns0="http://perf.getperf.com" xmlns="http://rmi.java/xsd">
            <ns0:siteKey>{{ .siteKey }}</ns0:siteKey>
            <ns0:hostname>{{ .hostname }}</ns0:hostname>
            <ns0:timestamp>{{ .timestamp }}</ns0:timestamp>
        </ns0:downloadCertificate xmlns:ns0="http://perf.getperf.com">
    </soapenv:Body>
</soapenv:Envelope>
`,
  "checkAgent": `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:ns0="http://perf.getperf.com">
  <soapenv:Header>
    <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
    <wsa:To>{{ .ToURL }}</wsa:To>
    <wsa:Action>urn:checkAgent</wsa:Action>
  </soapenv:Header>
  <soapenv:Body>
    <checkAgent xmlns="http://perf.getperf.com">
      <siteKey>{{ .siteKey }}</siteKey>
      <hostname>{{ .hostname }}</hostname>
      <accessKey>{{ .accessKey }}</accessKey>
    </checkAgent>
  </soapenv:Body>
</soapenv:Envelope>
`,
  "getLatestVersion": `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:ns0="http://perf.getperf.com">
  <soapenv:Header>
    <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
    <wsa:To>{{ .ToURL }}</wsa:To>
    <wsa:Action>urn:getLatestVersion</wsa:Action>
  </soapenv:Header>
  <soapenv:Body>
    <getLatestVersion xmlns="http://perf.getperf.com">
    </getLatestVersion>
  </soapenv:Body>
</soapenv:Envelope>
`,
  "registAgent": `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:ns0="http://perf.getperf.com">
  <soapenv:Header>
    <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
    <wsa:To>{{ .ToURL }}</wsa:To>
    <wsa:Action>urn:registAgent</wsa:Action>
  </soapenv:Header>
  <soapenv:Body>
    <registAgent xmlns="http://perf.getperf.com">
      <siteKey>{{ .siteKey }}</siteKey>
      <hostname>{{ .hostname }}</hostname>
      <accessKey>{{ .accessKey }}</accessKey>
    </registAgent>
  </soapenv:Body>
</soapenv:Envelope>

`,
  "sendMessage": `
<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/">
  <soap-env:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
    <wsa:Action>urn:sendMessage</wsa:Action>
    <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
    <wsa:To>{{ .ToURL }}</wsa:To>
  </soap-env:Header>
  <soap-env:Body>
    <ns0:sendMessage xmlns:ns0="http://perf.getperf.com">
      <ns0:siteKey>{{ .siteKey }}</ns0:siteKey>
      <ns0:hostname>{{ .hostname }}</ns0:hostname>
      <ns0:severity>{{ .severity }}</ns0:severity>
      <ns0:message>{{ .message }}</ns0:message>
    </ns0:sendMessage>
  </soap-env:Body>
</soap-env:Envelope>
`,
  "sendData": `
    <soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/">
    <soap-env:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
      <wsa:Action>urn:sendData</wsa:Action>
      <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
      <wsa:To>{{ .ToURL }}</wsa:To>
    </soap-env:Header>
    <soap-env:Body>
      <ns0:sendData xmlns:ns0="http://perf.getperf.com">
        <ns0:siteKey>{{ .siteKey }}</ns0:siteKey>
        <ns0:filename>{{ .filename }}</ns0:filename>
      </ns0:sendData>
    </soap-env:Body>
  </soap-env:Envelope>
`,
  "reserveSender": `
<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/">
  <soap-env:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
    <wsa:Action>urn:reserveSender</wsa:Action>
    <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
    <wsa:To>{{ .ToURL }}</wsa:To>
  </soap-env:Header>
  <soap-env:Body>
    <ns0:reserveSender xmlns:ns0="http://perf.getperf.com">
      <ns0:siteKey>{{ .siteKey }}</ns0:siteKey>
      <ns0:filename>{{ .filename }}</ns0:filename>
      <ns0:size>{{ .size }}</ns0:size>
    </ns0:reserveSender>
  </soap-env:Body>
</soap-env:Envelope>
`,
  "downloadUpdateModule": `
<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/">
  <soap-env:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
    <wsa:Action>urn:downloadUpdateModule</wsa:Action>
    <wsa:MessageID>urn:uuid:{{ .MessageID }}</wsa:MessageID>
    <wsa:To>{{ .ToURL }}</wsa:To>
  </soap-env:Header>
  <soap-env:Body>
    <ns0:downloadUpdateModule xmlns:ns0="http://perf.getperf.com">
      <ns0:moduleTag>{{ .moduleTag }}</ns0:moduleTag>
      <ns0:majorVer>{{ .majorVer }}</ns0:majorVer>
      <ns0:build>{{ .build }}</ns0:build>
    </ns0:downloadUpdateModule>
  </soap-env:Body>
</soap-env:Envelope>
`,
}
