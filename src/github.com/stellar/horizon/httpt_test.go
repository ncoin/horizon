package horizon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stellar/horizon/test"
)

type HTTPT struct {
	Assert *Assertions
	App    *App
	RH     test.RequestHelper
	*test.T
}

// StartHTTPTest is a helper function to setup a new test that will make http
// requests. Pair it with a deferred call to FinishHTTPTest.
func StartHTTPTest(t *testing.T, scenario string) *HTTPT {
	ret := &HTTPT{T: test.Start(t).Scenario(scenario)}
	ret.App = NewTestApp()
	ret.RH = test.NewRequestHelper(ret.App.web.router)
	ret.Assert = &Assertions{ret.T.Assert}

	return ret
}

func (ht *HTTPT) Get(
	path string,
	fn ...func(*http.Request),
) *httptest.ResponseRecorder {
	return ht.RH.Get(path, fn...)
}

func (ht *HTTPT) Finish() {
	ht.T.Finish()
	ht.App.Close()
}

// ReapHistory causes the test server to run `DeleteUnretainedHistory`, after
// setting the retention count to the provided number.
func (ht *HTTPT) ReapHistory(retention uint) {
	ht.App.reaper.RetentionCount = retention
	err := ht.App.DeleteUnretainedHistory()
	ht.Require.NoError(err)
	ht.App.UpdateLedgerState()
}