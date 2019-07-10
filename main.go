package main

import(
  "fmt"
  "os"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "time"
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

type Activities []struct {
  Name string `json:"name"`
  Distance float32 `json:"distance"`
  Date string `json:"start_date"`
}

func launchHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech("You have successfully launched a new session.")
	echoResponse.EndSession(false)
}

func sessionEndedHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
	echoResponse.OutputSpeech("Session ended.")
}

func intentHandler(request *skillserver.EchoRequest, echoResponse *skillserver.EchoResponse) {
  var response *skillserver.EchoResponse

  switch request.GetIntentName() {
  case "GetWeekDistance":
    response = handleGetWeekDistanceIntent(request)
  default:
    echoResponse.OutputSpeech(fmt.Sprintf("You have invoked the %s intent.", request.GetIntentName()))
  }

  if response == nil {
    response = skillserver.NewEchoResponse()
    response.OutputSpeech("Sorry, something went wrong")
  }

  *echoResponse = *response
}

func handleGetWeekDistanceIntent(request *skillserver.EchoRequest) *skillserver.EchoResponse {
  var bearer = "Bearer " + request.Session.User.AccessToken

  url := "https://www.strava.com/api/v3/athlete/activities?after=1562554008"

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatalln(err)
  }

  req.Header.Add("Authorization", bearer)

  client := &http.Client{Timeout: 10 * time.Second}
  res, err := client.Do(req)
  if err != nil {
    log.Fatalln(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    log.Fatalln(err)
  }

  var responseObject Activities
  json.Unmarshal(body, &responseObject)
  log.Println(len(responseObject))

  var totalDistance float32

  for i := 0; i < len(responseObject); i++ {
    log.Println(responseObject[i].Distance)
    totalDistance += responseObject[i].Distance
  }

  response := skillserver.NewEchoResponse()
  response.OutputSpeech(fmt.Sprintf("You have run %.2f kilometers this week.", totalDistance / 1000))

  return response
}

func main() {
  skillserver.Run(applications, "8081")
}
