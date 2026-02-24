# EZ Baccarat Simulator (Go)

*Read this in other languages: [English](README.md), [简体中文](README_zh.md)*

A Monte Carlo simulation engine and interactive command-line game for **EZ Baccarat** (免佣百家乐), written in Go. 

This project mathematically implements the strict drawing rules of EZ Baccarat, including the official payouts and statistical probabilities for the popular **Dragon 7** and **Panda 8** side bets.

## Features
* **Mathematically Accurate Engine**: Strictly follows the [California Bureau of Gambling Control (BGC)](https://oag.ca.gov/sites/all/files/agweb/pdfs/gambling/BGC_ez_baccarat_panda_8.pdf) third-card drawing matrix and Natural 8/9 rules.
* **Concurrent Monte Carlo Simulation**: Uses goroutines and channels to simulate millions of hands per second, outputting exact House Edge and probabilities that match theoretical Baccarat academic limits.
* **Interactive CLI Gameplay**: Play a standard game natively in your terminal.
* **Special Bet Support**: Full rule coverage for EZ Baccarat's core "Dragon 7" (1:40) and "Panda 8" (1:25) side bets.
* **Player Data Persistence**: Saves and tracks simulated player bankrolls and game history logs across sessions.

## The Rules of EZ Baccarat

Unlike traditional Baccarat, EZ Baccarat eliminates the 5% commission charged on winning Banker bets and replaces it with a specific winning condition.

### Core Payouts
* **Player Win (闲赢)**: Pays 1 to 1.
* **Banker Win (庄赢)**: Pays 1 to 1 (No Commission / 免水).
* **Tie (和局)**: Pays 8 to 1.

### The Exceptions (Side Bets)
1. **Dragon 7 (龙七)**:
   * Occurs when the Banker wins with exactly **3 cards** totaling **7 points**.
   * Any bet on the **Banker** becomes a **Push** (returns original bet, no win/loss). This replaces the house commission.
   * Any specific side bet placed on **Dragon 7** pays **40 to 1**.

2. **Panda 8 (熊猫8)**: 
   * Occurs when the Player wins with exactly **3 cards** totaling **8 points**.
   * Any bet on the **Player** pays 1 to 1 normally.
   * Any specific side bet placed on **Panda 8** pays **25 to 1**.

### Drawing Rules
The exact third-card hit/stand matrix relies on player's total and Banker's total. Most notably: **If either the Player or the Banker is dealt an 8 or 9 on the first two cards (a "**Natural**"), the hand is over.** Neither side may draw a third card.

## Simulation Results

Here is an example output from a 100,000,000 round Monte Carlo simulation:

```text
Starting simulation of 100000000 rounds using 30 workers...

=== Simulation Complete ===
Total Rounds: 100000000
Time Taken:   20.634266008s (4846308 rounds/sec)

Outcome              | Count        | Simulated %  | Expected %  
------------------------------------------------------------------
Player (Total)       |     44623956 |     44.6240% |     44.6247%
  ↳ Panda 8          |      3454935 |      3.4549% |      3.4543%
Banker (Non-Dragon)  |     43609570 |     43.6096% |     43.6064%
Tie                  |      9513636 |      9.5136% |      9.5156%
Dragon 7             |      2252838 |      2.2528% |      2.2534%
------------------------------------------------------------------
Total                |    100000000 |    100.0000% |    100.0000%
==================================================================

Bet Type ($1/hand)   | Net Profit ($)   | Simulated EV    | Expected EV    
-----------------------------------------------------------------------
Banker               |         -1014386 |        -1.0144% |        -1.0183%
Player               |         -1238452 |        -1.2385% |        -1.2351%
Tie                  |        -14377276 |       -14.3773% |       -14.3596%
Dragon 7             |         -7633642 |        -7.6336% |        -7.6106%
Panda 8              |        -10171690 |       -10.1717% |       -10.1882%
=======================================================================
```

*(Detailed theoretical combinations vs expected values formulas can be found in the `ez_baccarat_requirements.md` specifications)*

## Usage

### 1. Build the Project
```bash
go build -o ez_baccarat .
```

### 2. Run Interactive CLI Mode
Start an interactive terminal session where you can place bets with fake currency as a simulated player:
```bash
# Start a game as a default player
./ez_baccarat

# Or create a new user profile with starting bankroll
./ez_baccarat --create_player --player="Alice" --initial_balance=10000

# Continue playing as an existing user
./ez_baccarat --player="Alice"
```

### 3. Run Monte Carlo Simulation Mode
Run a multi-threaded headless probability simulation to calculate output occurrences and mathematical edge:

```bash
# Run 1,000,000 hands across 8 CPU threads
./ez_baccarat --simulate=1000000 --workers=8
```
