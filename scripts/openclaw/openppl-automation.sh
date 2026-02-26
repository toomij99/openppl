#!/usr/bin/env bash
set -euo pipefail

OPENPPL_BIN="${OPENPPL_BIN:-openppl}"

die() {
  printf '{"version":"v1","result_state":"rejected","error":{"code":"wrapper.invalid_arguments","message":"%s"}}\n' "$1" >&2
  exit 2
}

if [[ $# -lt 1 ]]; then
  die "expected subcommand: status | remind"
fi

subcommand="$1"
shift

case "$subcommand" in
  status)
    if [[ $# -ne 0 ]]; then
      die "status accepts no extra arguments"
    fi
    exec "$OPENPPL_BIN" automation status
    ;;
  remind)
    request_id=""
    actor_scope="telegram:default"
    while [[ $# -gt 0 ]]; do
      case "$1" in
        --request-id)
          [[ $# -ge 2 ]] || die "--request-id requires a value"
          request_id="$2"
          shift 2
          ;;
        --actor-scope)
          [[ $# -ge 2 ]] || die "--actor-scope requires a value"
          actor_scope="$2"
          shift 2
          ;;
        *)
          die "unsupported argument: $1"
          ;;
      esac
    done

    [[ -n "$request_id" ]] || die "--request-id is required"
    exec "$OPENPPL_BIN" automation action --name remind --request-id "$request_id" --actor-scope "$actor_scope"
    ;;
  *)
    die "unsupported subcommand: $subcommand"
    ;;
esac
