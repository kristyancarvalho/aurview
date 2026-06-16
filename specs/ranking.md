# Ranking

Ranking is isolated in `internal/ranking`.

Signals:

- exact name match
- prefix match
- substring match
- description match
- votes
- popularity
- recency / last update
- package not flagged out-of-date
- maintained vs orphaned package

Results from multiple sources are ranked together. Duplicate package names from different sources remain separate rows.
