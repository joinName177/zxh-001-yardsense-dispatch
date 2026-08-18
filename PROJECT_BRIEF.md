# YardSense Dispatch product brief

YardSense Dispatch is an operations console for a logistics yard where safety, loading, and equipment incidents must remain visible while teams rotate shifts. A dispatcher records a ticket with zone, severity, due time, and tags; only active workers who cover the zone can receive it. Every write is persisted in a local JSON document so a restarted console retains the board, assignments, notes, revision numbers, and audit trail.

The core flows are ticket intake, staffing management, guarded assignment, status transition, note capture, dashboard aggregation, and restoration after reopening the data file. The static frontend exposes a live dispatch board while the HTTP API supports integrations.
