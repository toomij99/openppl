#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

journal="$tmp_dir/journal.txt"
fake_bin="$tmp_dir/openppl"

cat > "$fake_bin" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

journal_file="${OPENPPL_AUTOMATION_JOURNAL:?OPENPPL_AUTOMATION_JOURNAL not set}"

if [[ "$1" != "automation" ]]; then
  echo '{"version":"v1","result_state":"error","error":{"code":"invalid.command","message":"expected automation"}}' >&2
  exit 1
fi

shift
case "$1" in
  status)
    echo '{"version":"v1","result_state":"ok","timestamp":"2026-01-01T00:00:00Z","status":{"summary":{"total_tasks":1,"completed_tasks":0,"pending_tasks":1},"next_tasks":[{"date":"2026-01-02","category":"Theory","title":"Task"}]}}'
    ;;
  action)
    shift
    name=""
    request_id=""
    actor_scope="default"
    while [[ $# -gt 0 ]]; do
      case "$1" in
        --name)
          name="$2"
          shift 2
          ;;
        --request-id)
          request_id="$2"
          shift 2
          ;;
        --actor-scope)
          actor_scope="$2"
          shift 2
          ;;
        *)
          shift
          ;;
      esac
    done

    key="${name}:${request_id}:${actor_scope}"
    if grep -qx "$key" "$journal_file" 2>/dev/null; then
      echo '{"version":"v1","result_state":"replayed","action":{"action_name":"remind","created_count":1}}'
      exit 0
    fi

    echo "$key" >> "$journal_file"
    echo '{"version":"v1","result_state":"executed","action":{"action_name":"remind","created_count":1}}'
    ;;
  *)
    echo '{"version":"v1","result_state":"error","error":{"code":"invalid.subcommand","message":"unsupported"}}' >&2
    exit 1
    ;;
esac
EOF

chmod +x "$fake_bin"
chmod +x "$repo_root/scripts/openclaw/openppl-automation.sh"
chmod +x "$repo_root/scripts/openclaw/deploy.sh"
chmod +x "$repo_root/scripts/openclaw/smoke.sh"

OPENPPL_BIN="$fake_bin" \
OPENPPL_AUTOMATION_JOURNAL="$journal" \
OPENPPL_AUTOMATION_WRAPPER="$repo_root/scripts/openclaw/openppl-automation.sh" \
OPENPPL_SMOKE_REQUEST_ID="test-request-1" \
OPENPPL_SMOKE_ACTOR_SCOPE="telegram:test" \
  bash "$repo_root/scripts/openclaw/smoke.sh"

OPENPPL_BIN="$fake_bin" \
OPENPPL_AUTOMATION_JOURNAL="$journal" \
  bash "$repo_root/scripts/openclaw/deploy.sh" dry-run

echo "openclaw smoke test passed"
