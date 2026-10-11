## Permanent Agent Handoff Rule
Every coding agent working in this repository must maintain:
`context/AGENT_HANDOFF.md`
Before ending a work session, update that file with:
- Task performed
- Files changed
- Tests/commands executed with PASS/FAIL
- Errors and unresolved problems
- `git status --short`
- `git diff --stat`
- Remaining work
- `TASK_COMPLETE=YES/NO`
- `READY_TO_COMMIT=YES/NO`
Rules:
- Keep it concise.
- Record verified facts only.
- Never report PASS unless the command actually ran successfully.
- Never hide failed commands.
- Use `NOT_VERIFIED` when something could not be tested.
- Update the existing handoff instead of creating new handoff files.
- The handoff must reflect the repository state at the end of the agent's work.

## Terminal Evidence Rule

For every important verification command, preserve the complete stdout and stderr in:

`context/TERMINAL_SESSION.log`

Before running verification pipelines, enable pipeline failure propagation:

`set -o pipefail`

Then capture command output with:

`COMMAND 2>&1 | tee -a context/TERMINAL_SESSION.log`

Logging must preserve the original command's failure status. A failed command must never be reported as PASS because `tee` succeeded.

At the start of each work session, append a timestamped session separator to `context/TERMINAL_SESSION.log`. The log is append-only during that work session. Rotate or truncate the log between sessions when appropriate so it does not grow without bound.

Always capture build, test, race, vet, lint, security scan, formatting,
and Git verification commands when their output is safe to persist.

Never log secrets, tokens, credentials, `.env` contents, private keys, authentication headers, or sensitive environment variables. If a command may expose sensitive data, do not capture its raw output; record only a sanitized result in `context/AGENT_HANDOFF.md`.

Update `context/AGENT_HANDOFF.md` with the concise result.

Never treat the handoff summary as a replacement for raw command evidence.
