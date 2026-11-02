import { useEffect, useState } from "react";
import { RefreshCw, Save, Trash2, Zap } from "lucide-react";
import { useLive } from "../api/context";
import { Empty, Pill, type PageProps } from "../ui";
import "../hot-update.css";

type HotUpdateData = {
  sell?: Record<string, number>;
  upgrade?: Record<string, number>;
  epicRate?: Record<string, number>;
  magicBox?: Record<string, { itemId: number; weight: number; count: number }[]>;
};

export function HotUpdatePage({ state, notify }: PageProps) {
  const live = useLive();
  const data = (live
    ? (state.api?.resources.hotUpdate as HotUpdateData | undefined)
    : undefined) as HotUpdateData | undefined;
  const canRead = !live || !!state.api?.features["hot-update.read"];
  const canWrite = !live || !!state.api?.features["hot-update.write"];
  const busy = !!state.api?.busy;
  const [q, setQ] = useState("");

  useEffect(() => {
    if (!live) return;
    void live.guard(() => live.hotUpdate());
    return () => live.cancelScope("hot-update");
  }, [live]);

  // ── 出售价 ──
  const [sellItem, setSellItem] = useState("");
  const [sellGold, setSellGold] = useState("");
  // ── 强化/增幅 ──
  const [upgradeMode, setUpgradeMode] = useState("强化");
  const [upgradeLevel, setUpgradeLevel] = useState("");
  const [upgradeWeight, setUpgradeWeight] = useState("");
  // ── 史诗掉率 ──
  const [epicCharacter, setEpicCharacter] = useState("");
  const [epicRate, setEpicRate] = useState("");
  // ── 魔盒 ──
  const [boxItem, setBoxItem] = useState("");
  const [boxGroups, setBoxGroups] = useState("");

  const sellRows = Object.entries(data?.sell || {})
    .filter(([k]) => k.includes(q.trim()))
    .map(([item, gold]) => ({ item, gold }));
  const upgradeRows = Object.entries(data?.upgrade || {})
    .filter(([k]) => k.includes(q.trim()))
    .map(([key, weight]) => {
      const [mode, level] = key.split(":");
      return { mode, level, weight };
    });
  const epicRows = Object.entries(data?.epicRate || {})
    .filter(([k]) => k.includes(q.trim()))
    .map(([character, rate]) => ({ character, rate }));
  const boxRows = Object.entries(data?.magicBox || {})
    .filter(([k]) => k.includes(q.trim()))
    .map(([box, groups]) => ({ box, groups }));

  async function run(action: () => Promise<unknown> | undefined, success: string) {
    try {
      await action();
      notify(success);
    } catch (e) {
      notify((e as Error).message);
    }
  }

  return (
    <div className="motion-section hot-update-page">
      <div className="hot-update-head glass">
        <div className="hot-update-title">
          <Zap size={19} />
          <div>
            <strong>游戏数值热更新</strong>
            <span>改库后 1 秒内对在线玩家生效，不重启游戏服务、不动 PVF 文件。</span>
          </div>
        </div>
        {live && (
          <button
            className="secondary"
            disabled={busy || !canRead}
            onClick={() => void live.guard(() => live.hotUpdate())}
          >
            <RefreshCw size={15} />
            刷新
          </button>
        )}
      </div>

      <div className="hot-update-toolbar">
        <input
          className="search"
          placeholder="搜索物品 / 角色 / 盒子编号"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
      </div>

      <div className="hot-update-grid">
        {/* 出售价 */}
        <section className="card hot-update-card">
          <header>
            <h3>物品出售价</h3>
            <Pill tone="purple">pvf_override</Pill>
          </header>
          <p className="muted">覆盖玩家把物品卖给 NPC 时回收的金币数。</p>
          <div className="hot-form">
            <label>
              物品 ID
              <input
                value={sellItem}
                onChange={(e) => setSellItem(e.target.value)}
                placeholder="例如 2600014"
                inputMode="numeric"
              />
            </label>
            <label>
              出售价（金币）
              <input
                value={sellGold}
                onChange={(e) => setSellGold(e.target.value)}
                placeholder="例如 100000"
                inputMode="numeric"
              />
            </label>
            <button
              className="primary"
              disabled={busy || !canWrite || !sellItem}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateSell({
                      itemId: Number(sellItem),
                      sellGold: Number(sellGold || 0),
                    }),
                  "出售价已设置，1 秒内生效",
                )
              }
            >
              <Save size={14} />
              设置
            </button>
            <button
              className="secondary"
              disabled={busy || !canWrite || !sellItem}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateSell({
                      itemId: Number(sellItem),
                      sellGold: 0,
                      clear: true,
                    }),
                  "出售价覆盖已清除",
                )
              }
            >
              <Trash2 size={14} />
              清除
            </button>
          </div>
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>物品 ID</th>
                  <th>出售价</th>
                  <th />
                </tr>
              </thead>
              <tbody>
                {sellRows.map((row) => (
                  <tr key={row.item}>
                    <td>{row.item}</td>
                    <td>{Number(row.gold).toLocaleString()}</td>
                    <td>
                      <button
                        className="text-button"
                        disabled={busy || !canWrite}
                        onClick={() =>
                          run(
                            () =>
                              live?.hotUpdateSell({
                                itemId: Number(row.item),
                                sellGold: 0,
                                clear: true,
                              }),
                            "已清除",
                          )
                        }
                      >
                        清除
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {!sellRows.length && (
            <Empty title="暂无出售价覆盖" detail="设置后显示在这里。" />
          )}
        </section>

        {/* 强化/增幅 */}
        <section className="card hot-update-card">
          <header>
            <h3>强化 / 增幅成功率</h3>
            <Pill tone="purple">upgrade_success_override</Pill>
          </header>
          <p className="muted">覆盖指定等级的成功率（万分比）。</p>
          <div className="hot-form">
            <label>
              模式
              <select
                value={upgradeMode}
                onChange={(e) => setUpgradeMode(e.target.value)}
              >
                <option value="强化">强化</option>
                <option value="增幅">增幅</option>
              </select>
            </label>
            <label>
              目标等级
              <input
                value={upgradeLevel}
                onChange={(e) => setUpgradeLevel(e.target.value)}
                placeholder="例如 12"
                inputMode="numeric"
              />
            </label>
            <label>
              成功率（万分比）
              <input
                value={upgradeWeight}
                onChange={(e) => setUpgradeWeight(e.target.value)}
                placeholder="例如 1000 = 10%"
                inputMode="numeric"
              />
            </label>
            <button
              className="primary"
              disabled={busy || !canWrite || !upgradeLevel}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateUpgrade({
                      mode: upgradeMode,
                      level: Number(upgradeLevel),
                      weight: Number(upgradeWeight || 0),
                    }),
                  "成功率已设置，1 秒内生效",
                )
              }
            >
              <Save size={14} />
              设置
            </button>
            <button
              className="secondary"
              disabled={busy || !canWrite || !upgradeLevel}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateUpgrade({
                      mode: upgradeMode,
                      level: Number(upgradeLevel),
                      weight: 0,
                      clear: true,
                    }),
                  "成功率覆盖已清除",
                )
              }
            >
              <Trash2 size={14} />
              清除
            </button>
          </div>
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>模式</th>
                  <th>等级</th>
                  <th>成功率</th>
                  <th />
                </tr>
              </thead>
              <tbody>
                {upgradeRows.map((row) => (
                  <tr key={`${row.mode}:${row.level}`}>
                    <td>{row.mode}</td>
                    <td>+{row.level}</td>
                    <td>{Number(row.weight).toLocaleString()} / 10000</td>
                    <td>
                      <button
                        className="text-button"
                        disabled={busy || !canWrite}
                        onClick={() =>
                          run(
                            () =>
                              live?.hotUpdateUpgrade({
                                mode: row.mode,
                                level: Number(row.level),
                                weight: 0,
                                clear: true,
                              }),
                            "已清除",
                          )
                        }
                      >
                        清除
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {!upgradeRows.length && (
            <Empty title="暂无成功率覆盖" detail="设置后显示在这里。" />
          )}
        </section>

        {/* 史诗掉率 */}
        <section className="card hot-update-card">
          <header>
            <h3>角色史诗掉率</h3>
            <Pill tone="purple">epic_rate_override</Pill>
          </header>
          <p className="muted">按角色覆盖史诗掉率倍率，1.0 = 原倍率。</p>
          <div className="hot-form">
            <label>
              角色 ID
              <input
                value={epicCharacter}
                onChange={(e) => setEpicCharacter(e.target.value)}
                placeholder="例如 77"
                inputMode="numeric"
              />
            </label>
            <label>
              倍率
              <input
                value={epicRate}
                onChange={(e) => setEpicRate(e.target.value)}
                placeholder="例如 2.0"
              />
            </label>
            <button
              className="primary"
              disabled={busy || !canWrite || !epicCharacter}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateEpic({
                      characterId: Number(epicCharacter),
                      rateMultiplier: Number(epicRate || 1),
                    }),
                  "史诗掉率已设置，1 秒内生效",
                )
              }
            >
              <Save size={14} />
              设置
            </button>
            <button
              className="secondary"
              disabled={busy || !canWrite || !epicCharacter}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateEpic({
                      characterId: Number(epicCharacter),
                      rateMultiplier: 1,
                      clear: true,
                    }),
                  "史诗掉率覆盖已清除",
                )
              }
            >
              <Trash2 size={14} />
              清除
            </button>
          </div>
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>角色 ID</th>
                  <th>倍率</th>
                  <th />
                </tr>
              </thead>
              <tbody>
                {epicRows.map((row) => (
                  <tr key={row.character}>
                    <td>{row.character}</td>
                    <td>{Number(row.rate).toFixed(2)}×</td>
                    <td>
                      <button
                        className="text-button"
                        disabled={busy || !canWrite}
                        onClick={() =>
                          run(
                            () =>
                              live?.hotUpdateEpic({
                                characterId: Number(row.character),
                                rateMultiplier: 1,
                                clear: true,
                              }),
                            "已清除",
                          )
                        }
                      >
                        清除
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {!epicRows.length && (
            <Empty title="暂无史诗掉率覆盖" detail="设置后显示在这里。" />
          )}
        </section>

        {/* 魔盒 */}
        <section className="card hot-update-card">
          <header>
            <h3>魔盒奖励</h3>
            <Pill tone="purple">magic_box_override</Pill>
          </header>
          <p className="muted">
            按盒子整组替换奖励清单与权重（赛丽亚的幸运 = 2682272）。
          </p>
          <div className="hot-form">
            <label>
              盒子物品 ID
              <input
                value={boxItem}
                onChange={(e) => setBoxItem(e.target.value)}
                placeholder="例如 2682272"
                inputMode="numeric"
              />
            </label>
            <label className="wide">
              奖励组（JSON）
              <textarea
                rows={4}
                value={boxGroups}
                onChange={(e) => setBoxGroups(e.target.value)}
                placeholder={'[{"itemId":2600014,"weight":60,"count":2}]'}
              />
            </label>
            <button
              className="primary"
              disabled={busy || !canWrite || !boxItem}
              onClick={() => {
                let groups: { itemId: number; weight: number; count: number }[];
                try {
                  groups = JSON.parse(boxGroups || "[]");
                } catch {
                  notify("奖励组 JSON 格式无效");
                  return;
                }
                void run(
                  () =>
                    live?.hotUpdateMagicBox({
                      boxItemId: Number(boxItem),
                      groups,
                    }),
                  "魔盒奖励已设置，1 秒内生效",
                );
              }}
            >
              <Save size={14} />
              设置
            </button>
            <button
              className="secondary"
              disabled={busy || !canWrite || !boxItem}
              onClick={() =>
                run(
                  () =>
                    live?.hotUpdateMagicBox({
                      boxItemId: Number(boxItem),
                      groups: [],
                      clear: true,
                    }),
                  "魔盒覆盖已清除",
                )
              }
            >
              <Trash2 size={14} />
              清除
            </button>
          </div>
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>盒子 ID</th>
                  <th>奖励条目</th>
                  <th />
                </tr>
              </thead>
              <tbody>
                {boxRows.map((row) => (
                  <tr key={row.box}>
                    <td>{row.box}</td>
                    <td className="muted">
                      {row.groups
                        .map(
                          (g) =>
                            `#${g.itemId}×${g.count}（${g.weight}%）`,
                        )
                        .join("、")}
                    </td>
                    <td>
                      <button
                        className="text-button"
                        disabled={busy || !canWrite}
                        onClick={() =>
                          run(
                            () =>
                              live?.hotUpdateMagicBox({
                                boxItemId: Number(row.box),
                                groups: [],
                                clear: true,
                              }),
                            "已清除",
                          )
                        }
                      >
                        清除
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {!boxRows.length && (
            <Empty title="暂无魔盒覆盖" detail="设置后显示在这里。" />
          )}
        </section>
      </div>
    </div>
  );
}
