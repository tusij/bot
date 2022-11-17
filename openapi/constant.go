package openapi

import "fmt"

const (
	domain = "api.sgroup.qq.com"

	scheme = "https"

	gatewayBotURI = "/gateway/bot"
	messagesURI   = "/channels/{channel_id}/messages"
)

func getURL(endpoint string) string {
	return fmt.Sprintf("%s://%s%s", scheme, domain, endpoint)
}
