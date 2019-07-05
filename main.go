package main

import(
  "fmt"
  "os"
  "github.com/mikeflynn/go-alexa/skillserver"
)

var (
  runningStatsAppID = os.Getenv("RUNNING_STATS_APP_ID")
	applications    = map[string]interface{}{
		"/echo/running-stats": skillserver.EchoApplication{
			AppID:          runningStatsAppID,
			OnIntent:       intentHandler,
			OnLaunch:       launchHandler,
			OnSessionEnded: sessionEndedHandler,
		},
	}
)

func launchHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech("You have successfully launched a new session.")
	echoResponse.EndSession(false)
}

func sessionEndedHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech("Session ended.")
}

func intentHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech(fmt.Sprintf("You have invoked the %s intent.", request.GetIntentName()))
}

func main() {
  skillserver.Run(applications, "8081")
}
