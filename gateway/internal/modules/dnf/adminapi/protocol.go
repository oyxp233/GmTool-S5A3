// Package adminapi defines the private game-runtime management contract.
// It contains no web panel, database connection, scheduler or game transport.
package adminapi

import (
	"dnf90.local/gateway/internal/modules/dnf/moderation"
	"dnf90.local/gateway/internal/modules/dnf/onlineevent"
	"dnf90.local/gateway/internal/modules/dnf/qqbinding"
	"dnf90.local/gateway/internal/modules/dnf/repository"
)

const Version = 1
const Path = "/api/game/admin"

type Request struct {
	PasswordReset   *PasswordReset              `json:"passwordReset,omitempty"`
	QQAck           *moderation.QQAck           `json:"qqAck,omitempty"`
	QQ              string                      `json:"qq,omitempty"`
	QQSettings      *qqbinding.Settings         `json:"qqSettings,omitempty"`
	QQComplete      *qqbinding.Complete         `json:"qqComplete,omitempty"`
	QQQuery         *qqbinding.Query            `json:"qqQuery,omitempty"`
	ModerationQuery *moderation.Query           `json:"moderationQuery,omitempty"`
	BanChange       *moderation.BanChange       `json:"banChange,omitempty"`
	AccountID       string                      `json:"accountId,omitempty"`
	Version         int                         `json:"version"`
	Action          string                      `json:"action"`
	CharacterID     uint16                      `json:"characterId,omitempty"`
	Revision        uint64                      `json:"revision,omitempty"`
	Changes         []onlineevent.ControlChange `json:"changes,omitempty"`
	ListType        byte                        `json:"listType,omitempty"`
	Slot            int16                       `json:"slot,omitempty"`
	Stack           *repository.ItemStack       `json:"stack,omitempty"`

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
	// Group 相同的条目组成一组，每次开盒每组按权重抽一条；省略即全部在第 0 组。
	Group int `json:"group,omitempty"`
}

type Response struct {
	PasswordReset      *PasswordResetResult         `json:"passwordReset,omitempty"`
	ErrorCode          string                       `json:"errorCode,omitempty"`
	BanTargets         *moderation.BanTargets       `json:"banTargets,omitempty"`
	QQActions          []moderation.QQAction        `json:"qqActions,omitempty"`
	QQBlocked          bool                         `json:"qqBlocked,omitempty"`
	QQSettings         *qqbinding.Settings          `json:"qqSettings,omitempty"`
	QQBinding          *qqbinding.Binding           `json:"qqBinding,omitempty"`
	QQBindings         *qqbinding.Page              `json:"qqBindings,omitempty"`
	Clears             *moderation.ClearPage        `json:"clears,omitempty"`
	Bans               *moderation.BanPage          `json:"bans,omitempty"`
	Ban                *moderation.Ban              `json:"ban,omitempty"`
	ModerationAccounts *moderation.AccountPage      `json:"moderationAccounts,omitempty"`
	CharacterEdit      *CharacterEditStatus         `json:"characterEdit,omitempty"`
	AccountEdit        *AccountEditStatus           `json:"accountEdit,omitempty"`
	OnlineSnapshot     *OnlineSnapshot              `json:"onlineSnapshot,omitempty"`
	Version            int                          `json:"version"`
	OK                 bool                         `json:"ok"`
	Error              string                       `json:"error,omitempty"`
	Online             bool                         `json:"online"`
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
