package server

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"k8s.io/test-infra/prow/hook"
	"net/http"
)

type GitHubEventHandler interface {
	HandleEvent(eventType, eventGUID string, payload []byte) error
}

// Server implements http.Handler. It validates incoming GitHub webhooks and
// then dispatches them to the appropriate plugins.
type Server struct {
	GitHubEventHandler GitHubEventHandler
	HmacSecret         []byte
	Log                *logrus.Entry
}

// ServeHTTP validates an incoming webhook and puts it into the event channel.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(k8s-prow): Move webhook handling logic out of hook binary so that we don't have to import all
	eventType, eventGUID, payload, ok := hook.ValidateWebhook(w, r, s.HmacSecret)
	if !ok {
		return
	}

	fmt.Fprint(w, "Event received. Have a nice day.")

	if err := s.GitHubEventHandler.HandleEvent(eventType, eventGUID, payload); err != nil {
		s.Log.WithError(err).Error("Error parsing event.")
	}
}
