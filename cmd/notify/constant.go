package main

const (
	toastBaseURL              = `https://mail.google.com/mail/u/`
	pathLogFilename           = `notify.log`
	pathImageFilename         = `mail.png`
	pathConfigFilename        = `config.json`
	pathCredentialsFilename   = `credentials.json`
	pathTokenFilename         = `token.json`
	pathUnreadFilename        = `unread.json`
	gmailAPIUserID            = `me`
	powerShellStreamSeparator = "GMAIL-EOL"
	powerShellStreamNewline   = "\r\n"
	powerShellInitCommand     = `
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
`
	commandTemplateCreateTemplateID = `toastNotifyCreateTemplate`
	commandTemplateCreateTemplate   = `
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$template = @'
<toast activationType="protocol" launch="{{.Launch}}">
	<visual>
		<binding template="ToastGeneric">
			<image placement="appLogoOverride" src="{{.ImageSrc}}" />
			<text><![CDATA[{{.Title}}]]></text>
			<text><![CDATA[{{.Message1}}]]></text>
			<text><![CDATA[{{.Message2}}]]></text>
		</binding>
	</visual>
	<audio {{if .SilentSound}} silent="true" {{end}} src="ms-winsoundevent:Notification.Mail" />
</toast>
'@

$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
$toast.Group = '{{.Group}}'
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('{{.AppID}}').Show($toast)
`
	commandTemplateRemoveTemplateID = `toastNotifyRemoveTemplate`
	commandTemplateRemoveTemplate   = `[Windows.UI.Notifications.ToastNotificationManager]::History.RemoveGroup('{{.Group}}' , '{{.AppID}}')`
	commandTemplateCountTemplateID  = `toastNotifyCountTemplate`
	commandTemplateCountTemplate    = `[Windows.UI.Notifications.ToastNotificationManager]::History.GetHistory('{{.AppID}}') | Where-Object {$_.Group -eq '{{.Group}}'}`
)
