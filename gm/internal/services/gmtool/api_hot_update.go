package gmtool

import (
	"net/http"

	"github.com/Cec1c/DNF90GM/internal/modules/dnf/adminapi"
)

// v1HotUpdateReady 检查游戏运行时通道可用。
func (s *Server) v1HotUpdateReady(w http.ResponseWriter, r *http.Request) bool {
	if s.deps.Runtime == nil {
		v1Fail(w, r, 503, "hot_update_unavailable", "游戏服务未接入热更新通道", true)
		return false
	}
	return true
}

// v1HotUpdateList 列出全部热更新覆盖（出售价 / 强化成功 / 史诗掉率 / 魔盒）。
func (s *Server) v1HotUpdateList(w http.ResponseWriter, r *http.Request) {
	if !s.v1HotUpdateReady(w, r) {
		return
	}
	sell, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: "pvf_override_list"})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_read_failed", "读取出售价覆盖失败："+err.Error(), true)
		return
	}
	upgrade, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: "upgrade_success_override_list"})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_read_failed", "读取强化成功覆盖失败："+err.Error(), true)
		return
	}
	epic, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: "epic_rate_override_list"})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_read_failed", "读取史诗掉率覆盖失败："+err.Error(), true)
		return
	}
	magic, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: "magic_box_override_list"})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_read_failed", "读取魔盒覆盖失败："+err.Error(), true)
		return
	}
	v1Reply(w, r, 200, v1Envelope{Data: map[string]any{
		"sell":       sell.Overrides,
		"upgrade":    upgrade.SuccessOverrides,
		"epicRate":   epic.EpicRateOverrides,
		"magicBox":   magic.MagicBoxOverrides,
	}})
}

// v1HotUpdateSellSet 设置物品出售价覆盖。
func (s *Server) v1HotUpdateSellSet(w http.ResponseWriter, r *http.Request) {
	if !s.v1HotUpdateReady(w, r) {
		return
	}
	var body struct {
		ItemID   uint32 `json:"itemId"`
		SellGold int64  `json:"sellGold"`
		Clear    bool   `json:"clear"`
	}
	if !v1Decode(w, r, &body) {
		return
	}
	if body.ItemID == 0 {
		v1Fail(w, r, 400, "invalid_body", "itemId 必须大于 0", false)
		return
	}
	action := "pvf_override_set"
	if body.Clear {
		action = "pvf_override_clear"
	}
	res, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: action, ItemID: body.ItemID, SellGold: body.SellGold})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_write_failed", "写入出售价覆盖失败："+err.Error(), true)
		return
	}
	v1Reply(w, r, 200, v1Envelope{Data: map[string]any{"overrideCount": res.OverrideCount}})
}

// v1HotUpdateUpgradeSet 设置强化/增幅成功率覆盖。
func (s *Server) v1HotUpdateUpgradeSet(w http.ResponseWriter, r *http.Request) {
	if !s.v1HotUpdateReady(w, r) {
		return
	}
	var body struct {
		Mode   string `json:"mode"`
		Level  int    `json:"level"`
		Weight int    `json:"weight"`
		Clear  bool   `json:"clear"`
	}
	if !v1Decode(w, r, &body) {
		return
	}
	if body.Mode == "" || body.Level < 1 {
		v1Fail(w, r, 400, "invalid_body", "mode 与 level 必填", false)
		return
	}
	action := "upgrade_success_override_set"
	if body.Clear {
		action = "upgrade_success_override_clear"
	}
	res, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: action, Mode: body.Mode, Level: body.Level, Weight: body.Weight})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_write_failed", "写入强化成功覆盖失败："+err.Error(), true)
		return
	}
	v1Reply(w, r, 200, v1Envelope{Data: map[string]any{"successOverrideCount": res.SuccessOverrideCount}})
}

// v1HotUpdateEpicSet 设置角色史诗掉率覆盖。
func (s *Server) v1HotUpdateEpicSet(w http.ResponseWriter, r *http.Request) {
	if !s.v1HotUpdateReady(w, r) {
		return
	}
	var body struct {
		CharacterID    uint16  `json:"characterId"`
		RateMultiplier float64 `json:"rateMultiplier"`
		Clear          bool    `json:"clear"`
	}
	if !v1Decode(w, r, &body) {
		return
	}
	if body.CharacterID == 0 {
		v1Fail(w, r, 400, "invalid_body", "characterId 必须大于 0", false)
		return
	}
	action := "epic_rate_override_set"
	if body.Clear {
		action = "epic_rate_override_clear"
	}
	res, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: action, CharacterID: body.CharacterID, RateMultiplier: body.RateMultiplier})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_write_failed", "写入史诗掉率覆盖失败："+err.Error(), true)
		return
	}
	v1Reply(w, r, 200, v1Envelope{Data: map[string]any{"epicRateOverrideCount": res.EpicRateOverrideCount}})
}

// v1HotUpdateMagicBoxSet 设置魔盒奖励覆盖。
func (s *Server) v1HotUpdateMagicBoxSet(w http.ResponseWriter, r *http.Request) {
	if !s.v1HotUpdateReady(w, r) {
		return
	}
	var body struct {
		BoxItemID  uint32                             `json:"boxItemId"`
		Groups     []adminapi.MagicBoxRewardGroupJSON `json:"groups"`
		Clear      bool                               `json:"clear"`
	}
	if !v1Decode(w, r, &body) {
		return
	}
	if body.BoxItemID == 0 {
		v1Fail(w, r, 400, "invalid_body", "boxItemId 必须大于 0", false)
		return
	}
	action := "magic_box_override_set"
	if body.Clear {
		action = "magic_box_override_clear"
	}
	res, err := s.deps.Runtime.Call(r.Context(), adminapi.Request{Action: action, ItemID: body.BoxItemID, MagicBoxRewardGroups: body.Groups})
	if err != nil {
		v1Fail(w, r, 502, "hot_update_write_failed", "写入魔盒覆盖失败："+err.Error(), true)
		return
	}
	v1Reply(w, r, 200, v1Envelope{Data: map[string]any{"magicBoxOverrideCount": res.MagicBoxOverrideCount}})
}
