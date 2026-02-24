# EZ Baccarat Flutter + gRPC 跨平台架构设计

**日期**: 2026-02-23

## 1. 概述与目标 (Overview & Objectives)
本设计的核心目标是将当前基于 Go 语言构建的纯命令行 EZ Baccarat 模拟器，全面升级扩建为一个支持 Web、iOS 和 Android 的成熟跨平台多人在线博弈应用。

通过前端采用 **Flutter**，以及后端网络层引入 **gRPC / Protocol Buffers** 强类型通信协议，系统将能够实现稳定 60fps+ 的 3D 发牌与筹码动画，同时拥有安全严谨的类型边界与极高的并发扩展能力。

## 2. 核心架构选型 (Core Architecture Decisions)

### 2.1 通信协议: 单次 gRPC 请求 (gRPC Unary Calls)
*   **决策**: 我们将采用标准的单次 gRPC 请求 (Unary RPC) 而不是服务端流式传输 (Server Streaming)。
*   **依据**: Baccarat 发牌与补牌的数学计算在 Go 引擎中几乎是须臾完成的（微秒级）。后端在一个请求内一次性把整局的载荷报文（包括发出的所有 5-6 张牌、最终获胜方、玩家资金结算）全部算好并传回。
*   **动画外包 (Animation Delegation)**: 由于结果已经确定，Flutter 前端将接手所有后续“悬念”。负责解析这条 JSON/Protobuf 数据，并通过严格掌控时间差的动画引擎自行播放翻牌过程（例如庄家的第三张牌强制要求玩家等 2 秒才慢慢亮出），以此营造刺激感。这使得服务器完全成为无状态计算节点，无需挂靠长连接线程，彻底解放了后端的并发限制。

### 2.2 状态管理: 虚拟联机牌桌机制 (Multiplayer Tables)
*   **决策**: 弃用单机模式，实现一个自带容量的虚拟“大厅 (Lobby)”与“牌桌 (Table)”共享状态流。
*   **依据**: 为了还原真实赌场的社交共鸣体验。Go 引擎会在内存中常驻几个活跃的“牌桌（例如最大支持 7 人同桌）”。连入同一个桌子的不同端设备（网页/手机）将共享同一个虚拟鞋盒 (Shoe)。大家面对相同的荷官物理发牌，并互相能在界面上看到彼此抛向桌面的筹码盲注。满桌将锁定进入权限。

### 2.3 身份鉴权与数据持久层 (Auth & Data Layer)
*   **决策**: 采用演进式开发路线，首期直接挂载第三方 OAuth 联合登录接入完整数据层。
*   **阶段 1 (OAuth & JWT)**: 让用户能够通过 Google / Apple 进行一键无缝登入。Go 服务端验证由第三方签发的 OAuth 令牌，核实后自己签发专属于本系统的 JWT，建立登录会话。
*   **存储引擎 (Data Storage)**: 深度引入 **PostgreSQL** 作为存放玩家真实余额流水的关系型凭证库 (Source of Truth)；同步引入 **Redis** 集群来接管超高频读写的游戏大厅房间状态展示（例如全服各桌的连入状态），使得关系型数据库免受海量查询攻击。

## 3. 核心大纲：全局数据交互流 (High-Level Data Flow)

1.  **进场 (Join)**: Flutter 端发起单次请求 `JoinLobby()` -> 服务器返回当前所有存活的 `TableList` (通过 Redis 提速)。
2.  **落座 (Sit)**: Flutter 按下对应并发起 `JoinTable(tableID)` -> 服务器拦截鉴权并在游戏循环中分配空座号 (1-7号位置)。
3.  **盲注期 (Bet Window)**: Go 服务器进入 `State: BETTING_OPEN`，等待并倒计时。此时手机端各玩家疯狂搓动 UI 拖放筹码，频发 `PlaceBet(amount, type)` 告诉服务器各自押了谁。
4.  **发牌运算 (Dealing)**: 服务器倒计时归零，全局锁死所有下注并进入 `State: DEALING`。随后瞬间基于物理模型调用完整的底层数学规则，推导牌靴及输赢结果。
5.  **定局回传 (Resolution)**: 服务器瞬间把 `HandResult` 封装（包含：闲家最终牌面、庄家最终牌面、输赢定调及资金变化记录）一口气拍回给桌上的这 7 名连接玩家。
6.  **动画大戏 (Animation)**: 由 Flutter 前端发力，逐帧播放全桌筹码飞行以及 3D 纸牌翻转出分，随后进入下一局下注期。

## 4. 下一步实操路径计划 (Next Steps for Implementation)
1.  优先构建 `.proto` 原型约束：利用 Protocol Buffers 先将大厅、座位表、每轮下注握手和牌组输赢核心数据结构严格定义成 Schema 文件。
2.  在当前 Go 源码上外包一层 gRPC 控制器，确保其能够调用 `es-Baccarat/rules` 等包中已被我们反复推敲好的数学函数库。
3.  预置并部署 PostgreSQL 数据结构（含玩家表 User Profile 与资金进出记录表 Transaction Log）。
4.  新建 Flutter 脚手架项目，优先在本地拉起一套极简能够直接向后端发起 gRPC 拨号拿取一局明牌结果的联调测试 demo UI。
