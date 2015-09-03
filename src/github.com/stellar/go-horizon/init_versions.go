package horizon

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os/exec"
  "strings"
  "encoding/json"
)

func initStellarCoreVersion(app *App) {
  if app.config.StellarCoreUrl == "" {
    return
  }

  resp, err := http.Get(fmt.Sprint(app.config.StellarCoreUrl,"/info"))

  if err != nil {
    app.log.Panic(app.ctx, err)
  }

  defer resp.Body.Close()
  contents, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    app.log.Panic(app.ctx, err)
  }

  var responseJson map[string]*json.RawMessage
  err = json.Unmarshal(contents, &responseJson)

  var serverInfo map[string]string
  err = json.Unmarshal(*responseJson["info"], &serverInfo)
  app.coreVersion = serverInfo["build"]
}

func initHorizonVersion(app *App) {
  version, err := exec.Command("git", "describe", "--always", "--dirty", "--tags").Output()
  if err != nil {
    app.log.Panic(app.ctx, err)
  }
  app.horizonVersion = strings.TrimSpace(string(version))
}

func init() {
  appInit.Add("stellarCoreVersion", initStellarCoreVersion, "app-context", "log")
  appInit.Add("horizonVersion", initHorizonVersion, "app-context", "log")
}
