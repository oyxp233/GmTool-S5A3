// Package adminapi defines the private game-runtime management contract.
// It contains no web panel, database connection, scheduler or game transport.
package adminapi

import (
	"time"

	"github.com/Cec1c/DNF90GM/internal/modules/dnf/onlineevent"
	"github.com/Cec1c/DNF90GM/internal/modules/dnf/repository"
)

const Version = 1
const Path = "/api/game/admin"

// 编辑准入协议名。服务端会逐字比对这两个常量，名字对不上就会判定
// "编辑状态不可用"，所以不能改。
const (
	CharacterEditProtocol = "character-admission-v1"
	AccountEditProtocol   = "account-row-v1"
)

// MaxOnlineSnapshotMembers 与服务端校验上限保持一致。
const MaxOnlineSnapshotMembers = 5000

type Request struct {
	Version     int                         `json:"version"`
	Action      string                      `json:"action"`
	AccountID   string                      `json:"accountId,omitempty"`
	CharacterID uint16                      `json:"characterId,omitempty"`
	Revision    uint64                      `json:"revision,omitempty"`
	Changes     []onlineevent.ControlChange `json:"changes,omitempty"`
	ListType    byte                        `json:"listType,omitempty"`
	Slot        int16                       `json:"slot,omitempty"`
	Stack       *repository.ItemStack       `json:"stack,omitempty"`

	// ITEM-SELL-OVERRIDE:pvf_override_* action 的入参。
	ItemID   uint32 `json:"itemId,omitempty"`
	SellGold int64  `json:"sellGold,omitempty"`
	Clear    bool   `json:"clear,omitempty"`

	// UPGRADE-SUCCESS-OVERRIDE:upgrade_success_override_* action 的入参。
	Mode   string `json:"mode,omitempty"`
	Level  int    `json:"level,omitempty"`
	Weight int    `json:"weight,omitempty"`

	// EPIC-RATE-OVERRIDE:epic_rate_override_* action 复用 CharacterID（原 Request 字段）。
	RateMultiplier float64 `json:"rateMultiplier,omitempty"`

	// MAGIC-BOX-OVERRIDE:magic_box_override_set 的入参。ItemID 复用为盒子 item id，
	// MagicBoxRewardGroups 是覆盖后的整组 (item, weight, count) 条目。
	MagicBoxRewardGroups []MagicBoxRewardGroupJSON `json:"magicBoxRewardGroups,omitempty"`
}

// MagicBoxRewardGroupJSON 是 GM 侧传入/返回的一个覆盖奖励组（item/weight/count）。
type MagicBoxRewardGroupJSON struct {
	ItemID int64 `json:"itemId"`
	Weight int64 `json:"weight"`
	Count  int64 `json:"count"`
}

type Response struct {
	Version            int                          `json:"version"`
	OK                 bool                         `json:"ok"`
	Error              string                       `json:"error,omitempty"`
	Online             bool                         `json:"online"`
	CharacterEdit      *CharacterEditStatus         `json:"characterEdit,omitempty"`
	AccountEdit        *AccountEditStatus           `json:"accountEdit,omitempty"`
	OnlineSnapshot     *OnlineSnapshot              `json:"onlineSnapshot,omitempty"`
	Activities         *onlineevent.ControlSnapshot `json:"activities,omitempty"`
	PremiumTypes       []int64                      `json:"premiumTypes,omitempty"`
	CrystalCubeItemIDs [6]int64                     `json:"crystalCubeItemIds,omitempty"`
	StackRow           []byte                       `json:"stackRow,omitempty"`

	// ITEM-SELL-OVERRIDE:pvf_override_* action 的返回。
	OverrideCount int              `json:"overrideCount,omitempty"`
	Overrides     map[uint32]int64 `json:"overrides,omitempty"`

	// UPGRADE-SUCCESS-OVERRIDE:upgrade_success_override_* action 的返回。
	SuccessOverrideCount int            `json:"successOverrideCount,omitempty"`
	SuccessOverrides     map[string]int `json:"successOverrides,omitempty"`

	// EPIC-RATE-OVERRIDE:epic_rate_override_* action 的返回。
	EpicRateOverrideCount int                `json:"epicRateOverrideCount,omitempty"`
	EpicRateOverrides     map[string]float64 `json:"epicRateOverrides,omitempty"`

	// MAGIC-BOX-OVERRIDE:magic_box_override_* action 的返回。
	MagicBoxOverrideCount int                                  `json:"magicBoxOverrideCount,omitempty"`
	MagicBoxOverrides     map[string][]MagicBoxRewardGroupJSON `json:"magicBoxOverrides,omitempty"`
}

// CharacterEditStatus 回答"这个角色现在能不能被离线改档"。
type CharacterEditStatus struct {
	Protocol    string `json:"protocol"`
	AccountID   string `json:"accountId"`
	CharacterID uint16 `json:"characterId"`
	Busy        bool   `json:"busy"`
}

// AccountEditStatus 回答"这个账号现在能不能被离线改档"（账号级仓库/金库等）。
type AccountEditStatus struct {
	Protocol  string `json:"protocol"`
	AccountID string `json:"accountId"`
	Busy      bool   `json:"busy"`
}

type OnlineMember struct {
	CharacterID uint16 `json:"characterId"`
	Generation  uint64 `json:"generation,string"`
}

// OnlineSnapshot 是"本进程活会话"的完整时间点视图，只含角色号与代次，
// 不含连接地址、凭据或任何玩家资产。
type OnlineSnapshot struct {
	ID         string         `json:"id"`
	ObservedAt time.Time      `json:"observedAt"`
	Complete   bool           `json:"complete"`
	Members    []OnlineMember `json:"members"`
}
