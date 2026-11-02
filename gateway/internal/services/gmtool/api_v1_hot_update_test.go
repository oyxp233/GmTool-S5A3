package gmtool

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dnf90.local/gateway/internal/modules/dnf/adminapi"
)

func hotUpdateV1Fixture(t *testing.T) (*Server, *httptest.Server, string) {
	t.Helper()
	s, _, token := activityV1Fixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+strings.Repeat("x", 32) {
			t.Error("missing private auth")
		}
		var in adminapi.Request
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			t.Error(err)
		}
		out := adminapi.Response{Version: 1, OK: true}
		switch in.Action {
		case "pvf_override_list":
			out.Overrides = map[uint32]int64{2600014: 100000}
		case "upgrade_success_override_list":
			out.SuccessOverrides = map[string]int{"强化:12": 1000}
		case "epic_rate_override_list":
			out.EpicRateOverrides = map[string]float64{"77": 2.0}
		case "magic_box_override_list":
			out.MagicBoxOverrides = map[string][]adminapi.MagicBoxRewardGroupJSON{"2682272": {{ItemID: 1001, Weight: 60, Count: 2}}}
		}
		_ = json.NewEncoder(w).Encode(out)
	}))
	client, err := adminapi.NewClient(server.URL+adminapi.Path, strings.Repeat("x", 32))
	if err != nil {
		server.Close()
		t.Fatal(err)
	}
	s.deps.Runtime = client
	return s, server, token
}

func TestHotUpdateV1RequiresSessionAndCSRF(t *testing.T) {
	s, server, _ := hotUpdateV1Fixture(t)
	defer server.Close()
	body := `{"itemId":2600014,"sellGold":100000}`
	if code, _ := operationCall(s, "", "hot-update/sell", body, "", nil); code != 401 {
		t.Fatal("unauthenticated write accepted", code)
	}
}

func TestHotUpdateV1SellSetRoutesToGame(t *testing.T) {
	s, server, token := hotUpdateV1Fixture(t)
	defer server.Close()
	if code, out := operationCall(s, token, "hot-update/sell", `{"itemId":2600014,"sellGold":100000}`, "", nil); code != 200 {
		t.Fatal(code, out)
	}
	// verify list read works
	if code, out := v1TestCall(t, s, token, "GET", "/api/v1/hot-update"); code != 200 {
		t.Fatal(code, out)
	} else if data, ok := out["data"].(map[string]any); !ok || data["sell"] == nil {
		t.Fatal("hot-update list missing sell data", out)
	}
}

func TestHotUpdateV1Validation(t *testing.T) {
	s, server, token := hotUpdateV1Fixture(t)
	defer server.Close()
	for _, tc := range []struct {
		path, body string
	}{
		{"hot-update/sell", `{"itemId":0,"sellGold":1}`},
		{"hot-update/upgrade", `{"mode":"","level":1,"weight":100}`},
		{"hot-update/upgrade", `{"mode":"强化","level":0,"weight":100}`},
		{"hot-update/epic-rate", `{"characterId":0,"rateMultiplier":2}`},
		{"hot-update/magic-box", `{"boxItemId":0,"groups":[]}`},
	} {
		if code, out := operationCall(s, token, tc.path, tc.body, "", nil); code != 400 {
			t.Fatalf("invalid accepted %s: %d %v", tc.path, code, out)
		}
	}
}
