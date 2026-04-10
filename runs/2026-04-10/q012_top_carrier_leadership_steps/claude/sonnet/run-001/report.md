# Yearly Carrier Leadership by Completed Flights (1987–2025)

## Overview

Across 39 years of BTS on-time data, three carriers have held the top position by completed flights, producing four distinct leadership eras separated by three transitions.

## Leadership Eras

| Era | Carrier | Years | Duration |
|-----|---------|-------|----------|
| 1 | DL (Delta) | 1987–1989 | 3 years |
| 2 | US (US Airways) | 1990–1991 | 2 years |
| 3 | DL (Delta) | 1992–1999 | 8 years |
| 4 | WN (Southwest) | 2000–2025 | 26 years (ongoing) |

Delta led first, lost to US Airways briefly, reclaimed the top spot, then ceded it to Southwest in 2000 — a position Southwest has held uninterrupted for 26 years.

## Leadership Transitions

There were exactly **three transition years** (leadership_changed = 1):

| Year | From | To | Gap at Transition (pct pts) | Gap at Transition (flights) |
|------|------|----|----------------------------|-----------------------------|
| 1990 | DL | US | 3.34 | 174,169 |
| 1992 | US | DL | 0.85 | 42,776 |
| 2000 | DL | WN | 0.48 | 26,200 |

**Sharpest transition: 1990.** US Airways seized the lead from Delta by the widest margin of any transition year — a 174,169-flight gap representing 3.34 percentage points. The 2000 handoff to Southwest was the quietest: WN edged DL by just 26,200 flights (0.48 pct pts), the narrowest winning margin of all three transitions.

## Key Metrics

**Longest reign:** Southwest (WN), 2000–2025 — **26 consecutive years** with no change in leadership. By 2025, WN had completed 1,262,665 flights at a 19.96% share.

**Closest race (smallest gap in any year):** 2000, the very year Southwest took over. WN led DL by just 0.48 percentage points (16.43% vs 15.95%), a margin of 26,200 flights across a combined ~5.5 million completed flights.

**Peak dominance (widest leader–runner gap):** 2009. Despite the post-financial-crisis contraction, Southwest held a **9.14 pct pts** lead over AA (17.67% vs 8.53%), the largest gap on record, driven as much by American's shrinkage (542,349 flights) as by WN's own volume.

**Peak leader share:** Southwest in 2017 at **23.45%** of all completed US flights (1,311,398 flights), the highest single-year market share in the dataset.

**Most contested era:** 1990–1999. US Airways and Delta traded the top two spots repeatedly, with gaps frequently below 1 pct pt. In 1994, DL led US by just 0.52 pct pts (26,743 flights).

## Sample Row

First result row (1987 — partial year in the dataset):

| Year | Leader | Leader Flights | Leader Share% | Runner | Runner Flights | Runner Share% | Gap Flights | Gap Pct Pts | Prev Leader | Changed |
|------|--------|---------------|--------------|--------|---------------|--------------|-------------|-------------|-------------|---------|
| 1987 | DL | 183,717 | 14.22 | AA | 163,165 | 12.63 | 20,552 | 1.59 | DL | 0 |

Note: 1987 shows lower absolute flight counts (~183K for DL vs ~749K in 1988) consistent with the BTS dataset beginning mid-year in 1987.

## Data Quality Notes

- **Shares sum sensibly:** Leader + runner shares in any year are well below 100%, consistent with 5–15+ active carriers splitting the remainder.
- **Gap accuracy:** All `gap_flights = leader_flights − runner_flights` values verified against raw rows.
- **Transitions correctly detected:** `leadership_changed = 1` appears in exactly 1990, 1992, and 2000 — confirmed by inspecting `prev_leader` against `leader_carrier` for each row.
- **2020 anomaly:** OO (SkyWest, a regional operator) appears as runner-up with a 12.98% share, reflecting the pandemic collapse of major network carriers while regional schedules were partly maintained under CARES Act agreements.
