package agent

import "testing"

func TestSoapAgentGetLatestBuild(t *testing.T) {
	result, err := soapSender.GetLatestBuild("Linux", "2")
	if err != nil {
		t.Errorf("get latest build %s", err)
	}
	t.Log(result)
}

func TestSoapAgentRegistAgent(t *testing.T) {
	result, err := soapSender.RegistAgent("site1", "host1", "81e6011f1c0660a8062dbe4ade4e910d841d36c4", "/tmp/test1.zip")
	if err != nil {
		t.Errorf("regist agent %s", err)
	}
	t.Log(result)
}

func TestSoapAgentDownloadUpdateModule(t *testing.T) {
	result, err := soapSender.DownloadUpdateModule("Linux", "2", "1", "/tmp/test2.zip")
	if err != nil {
		t.Errorf("regist agent %s", err)
	}
	t.Log(result)
}
func TestSoapAgentCheckHostStatus(t *testing.T) {
	result, err := soapSender.CheckHostStatus("site1", "host1", "81e6011f1c0660a8062dbe4ade4e910d841d36c4")
	if err != nil {
		t.Errorf("regist agent %s", err)
	}
	t.Log(result)
}

func TestSoapAgentReserveSender(t *testing.T) {
	err := soapSenderData.ReserveSender("site1", "arc_host1__Linux_20230506_0800.zip")
	if err != nil {
		t.Errorf("reserve file sender %s", err)
	}
}

func TestSoapAgentReserveFileSender(t *testing.T) {
	err := soapSenderData.ReserveFileSender("on", 10)
	if err != nil {
		t.Errorf("reserve file sender %s", err)
	}
}

func TestSoapAgentSendData(t *testing.T) {
	result, err := soapSenderData.SendData("site1", "arc_host1__Linux_20230506_0800.zip", "../arc_host1__Linux_20230506_0800.zip")
	if err != nil {
		t.Errorf("send data %s", err)
	}
	t.Log(result)
}

func TestSoapAgentSendMessage(t *testing.T) {
	result, err := soapSenderData.SendMessage("site1","host1","1","this is a test")

	if err != nil {
		t.Errorf("send message %s", err)
	}
	t.Log(result)
}

func TestSoapAgentDownloadCertificate(t *testing.T) {
	result, err := soapSenderData.DownloadCertificate("site1","host1","0", "/tmp/test1.zip")

	if err != nil {
		t.Errorf("send message %s", err)
	}
	t.Log(result)
}
