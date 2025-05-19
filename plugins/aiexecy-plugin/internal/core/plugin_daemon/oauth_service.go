package plugin_daemon

import (
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/internal/utils/stream"
	"aiexec-plugin/pkg/entities/oauth_entities"
	"aiexec-plugin/pkg/entities/requests"
)

func OAuthGetAuthorizationURL(
	session *session_manager.Session,
	request *requests.RequestOAuthGetAuthorizationURL,
) (*stream.Stream[oauth_entities.OAuthGetAuthorizationURLResult], error) {
	return GenericInvokePlugin[requests.RequestOAuthGetAuthorizationURL, oauth_entities.OAuthGetAuthorizationURLResult](
		session,
		request,
		1,
	)
}

func OAuthGetCredentials(
	session *session_manager.Session,
	request *requests.RequestOAuthGetCredentials,
) (*stream.Stream[oauth_entities.OAuthGetCredentialsResult], error) {
	return GenericInvokePlugin[requests.RequestOAuthGetCredentials, oauth_entities.OAuthGetCredentialsResult](
		session,
		request,
		1,
	)
}
