# OpenPPL Automation Skill (Bounded)

This skill is intentionally constrained to two intents:

1. `status`
2. `send reminder`

Unknown intents are denied by default.

## Allowed Intent Mapping

- Intent: `status`
  - Command: `scripts/openclaw/openppl-automation.sh status`
- Intent: `send reminder`
  - Command: `scripts/openclaw/openppl-automation.sh remind --request-id <request_id> --actor-scope <actor_scope>`

## Normalization Rules

- `status` does not accept additional free-text arguments.
- `send reminder` requires `request_id`; if absent, reject.
- `actor_scope` defaults to `telegram:default` when omitted.

## Deny Rules

- Deny all commands outside explicit mappings above.
- Deny shell passthrough, raw `openppl` invocation, and arbitrary CLI flags.
- Deny unsupported action names.
