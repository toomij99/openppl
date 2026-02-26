#!/usr/bin/env bash
set -euo pipefail

wrapper="${OPENPPL_AUTOMATION_WRAPPER:-scripts/openclaw/openppl-automation.sh}"
request_id="${OPENPPL_SMOKE_REQUEST_ID:-smoke-request-001}"
actor_scope="${OPENPPL_SMOKE_ACTOR_SCOPE:-telegram:smoke}"

[[ -x "$wrapper" ]] || chmod +x "$wrapper"

status_json="$($wrapper status)"
echo "$status_json" | grep -q '"result_state":"ok"'

first="$($wrapper remind --request-id "$request_id" --actor-scope "$actor_scope")"
echo "$first" | grep -q '"result_state":"executed"'

second="$($wrapper remind --request-id "$request_id" --actor-scope "$actor_scope")"
echo "$second" | grep -q '"result_state":"replayed"'

echo "smoke checks passed"
