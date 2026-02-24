# EZ Baccarat Flutter + gRPC Architecture Design

**Date**: 2026-02-23

## 1. Overview & Objectives
The goal is to expand the current Go-based command-line EZ Baccarat simulator into a full-fledged, multi-platform multiplayer casino application supporting Web, iOS, and Android. 

By leveraging **Flutter** for the frontend and **gRPC / Protocol Buffers** for the backend communication, the system will achieve consistent 60fps+ 3D animations, strict type-safe networking, and high scalability.

## 2. Core Architecture Decisions

### 2.1 Communication Protocol: gRPC Unary Calls
*   **Decision**: We will use standard gRPC Unary calls (`Unary RPC`) rather than Server Streaming.
*   **Reasoning**: Baccarat math calculations (dealing and drawing the 3rd card) resolve instantaneously in the Go engine. The backend will compute the entire hand outcome in microseconds and return the complete payload (all cards drawn, winning hand, payoffs) in a single response.
*   **Animation Delegation**: The Flutter client is purely responsible for parsing this JSON/Protobuf payload and meticulously staging the delay animations (e.g., waiting 2 seconds before flipping the banker's 3rd card) to build suspense. The server remains stateless and avoids the overhead of managing long-lived streams per hand.

### 2.2 State Management: Multiplayer Tables (Lobbies)
*   **Decision**: Implement a virtual "Table" and "Lobby" architecture where multiple players can join and share the same game session.
*   **Reasoning**: Real-world casino experiences are deeply social. Go will maintain global memory states for several active "Tables" (e.g., up to 7 players per table). Players sharing a table will witness the exact same shoe, the exact same dealer outcomes, and see each other's bets land on the felt in real-time. Full tables will be locked to observers.

### 2.3 Authentication & Data Layer
*   **Decision**: Start with 3rd-party OAuth (C) while building the foundation for a robust relational database (B).
*   **Phase 1 (OAuth & JWT)**: Users will log in seamlessly using Google/Apple Sign-In. The Go server will validate the OAuth token, generate an internal JWT, and establish the session.
*   **Data Storage**: PostgreSQL will serve as the persistent source of truth for player balances and historical audit logs. Redis will be introduced to handle high-frequency, ephemeral state reading, such as broadcasting live lobby table capacities without hitting the SQL disk.

## 3. High-Level Data Flow

1.  **Join**: Flutter Client -> `JoinLobby()` -> Server returns active `TableList`.
2.  **Sit**: Flutter Client -> `JoinTable(tableID)` -> Server allocates a seat (1-7).
3.  **Bet Window**: Server broadcasts `State: BETTING_OPEN`. Clients drag and drop chips. Client calls `PlaceBet(amount, type)`.
4.  **Dealing**: Server broadcasts `State: DEALING`. It resolves the deterministic Baccarat algorithm.
5.  **Resolution**: Server replies with `HandResult` (Player Cards, Banker Cards, Winner, Balance Deltas).
6.  **Animation**: Flutter initiates 3D card flipping sequence. Chips animate toward winners. Back to Step 3.

## 4. Next Steps for Implementation
1.  Define `.proto` schema files for Lobby, Table, Core Hand, and Bet structures.
2.  Refactor the current `es-Baccarat/rules` and `es-Baccarat/model` to be imported by the new gRPC handler.
3.  Set up the PostgreSQL schema for player balances.
4.  Scaffold the initial Flutter UI focusing on the 3D card layout and single-player gRPC invocation.
