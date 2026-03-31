`review.md` written. Summary of the verdict:

**PASS** — all three sections are materially correct and evidence-supported.

Key checks that passed:
- All 7 required `main` fields present; 10 rows, all `hop_count=8`, `WN`, ordered by `flight_date DESC`
- Tail_Number NULL not filtered
- q1 lists all 10 routes individually; tier math 2+4+4=10
- q2 correctly scoped to the 10 main-route airports; 45 airports verified, all joined cleanly

Two minor prose imprecisions flagged in q2 (the "no geographic backtracking" claim fails for `BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX`, and the MDW-vs-ORD Southwest fingerprint is muddled because ORD appears in the data). Neither affects the numeric results or the section-level conclusions.
