package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	fileName := "fyne-test.ps1"
	tmpFilePath := filepath.Join(os.TempDir(), fileName)

	title := "Fyneの練習"
	content := "通知のテスト"
	script := notificationTemplate
	err := ioutil.WriteFile(tmpFilePath, []byte(script), 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := os.Remove(tmpFilePath); err != nil {
			log.Println(err)
		}
	}()

	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", "-File", tmpFilePath, "-title", title, "-content", content)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out), err)
	}
}

const notificationTemplate = `Param( $title, $content )

[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("appID").Show($toast);`
