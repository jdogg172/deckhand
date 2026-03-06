# OpenShift Local Notes

## Role in this plan
OpenShift Local is the local validation lane, not the first development lane.

Use it to verify:
- OpenShift projects
- Routes
- OpenShift Pipelines operator behavior
- OpenShift-specific API discovery and graceful fallback

## What Codex should implement for OpenShift awareness
- detect OpenShift API groups without hard failure
- treat projects/namespaces cleanly
- support Route listing when route API is available
- avoid assuming OpenShift resources always exist
- if Tekton/OpenShift Pipelines are missing, keep the TUI usable

## Practical local checks
- Start CRC and log in (`crc start`, `oc login ...`)
- Confirm project APIs: `oc api-resources | findstr /i project`
- Confirm route APIs: `oc api-resources | findstr /i route`
- Run Deckhand and validate:
	- project-aware namespace switching
	- Route listing when available
	- no panic / clear status if Route API is missing
	- Pipeline Mode still usable when Tekton APIs are missing
