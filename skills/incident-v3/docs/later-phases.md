# Later Phases — Phase 2 through Phase 4

Read this file when entering Phase 2. Do NOT preload.

---

## Phase 2: Decentralized Verification (Autonomous Agent Loop)

Each analyst-DA pair runs an independent Q&A loop. The DA drives the loop: receives findings, questions the analyst, requests scoring, and retries if needed. The orchestrator does NOT relay messages — it only collects final results.

### Architecture

```
analyst-1 ──findings──→ da-1 ──score request──→ shared-scorer
analyst-2 ──findings──→ da-2 ──score request──→ shared-scorer
analyst-N ──findings──→ da-N ──score request──→ shared-scorer

Each pair runs independently. Scorer processes requests FIFO.
DA reports verified findings to orchestrator (team-lead) when done.
```

### Agent Communication Protocol

| From | To | Message | When |
|------|-----|---------|------|
| Analyst | DA | Raw findings | After investigation complete |
| DA | Analyst | Socratic questions (2-4) | Each Q&A round |
| Analyst | DA | Clarification responses | After receiving questions |
| DA | Shared-Scorer | Score request (findings + QA history) | When DA thinks ambiguity is resolved |
| Shared-Scorer | DA | Score JSON (PASS/FAIL + improvement_hint) | After scoring |
| DA | Analyst | Follow-up questions | If scorer returns FAIL |
| DA | Orchestrator | Verified findings + final score | On PASS or FORCE PASS |

### Step 2.1: Wait for DA Reports

The orchestrator waits for all DAs to report via `SendMessage`. Each DA reports:
- **Verified findings** (post-Q&A)
- **Ambiguity score** (from scorer)
- **Q&A summary** (rounds completed, key clarifications)
- **Verdict**: PASS or FORCE PASS

As each DA reports, persist results:
- Write verified findings to `.omc/state/incident-{short-id}/verified-findings-{analyst-id}.md`
- Write score to `.omc/state/incident-{short-id}/ambiguity-{analyst-id}.json`

### Step 2.2: Compile Verified Findings

After ALL DAs have reported:

1. Compile all verified findings into `.omc/state/incident-{short-id}/analyst-findings.md`
2. Include ambiguity scores summary table
3. Flag any FORCE PASS analysts for user attention

### Phase 2 Exit Gate

- [ ] All DAs have reported verified findings
- [ ] All scores persisted
- [ ] Compiled findings written to `analyst-findings.md`

→ **NEXT ACTION: Proceed to Phase 2.5 — Tribunal Decision.**

---

## Phase 2.5: Tribunal Decision

Tribunal is now a user decision, not an automatic trigger.

### Step 2.5.1: Present Summary

Show the user a summary of verified findings:
- Per-analyst ambiguity scores
- Any FORCE PASS analysts (highlighted)
- Key findings overview

### Step 2.5.2: Ask User

```
AskUserQuestion(
  header: "Tribunal",
  question: "All analyst findings have been verified through Socratic Q&A. Would you like to send them to a Tribunal for additional review?",
  options: [
    "Skip — proceed to report",
    "Request Tribunal"
  ]
)
```

### Step 2.5.3: If Tribunal Requested

1. Compile findings package (~10-15K tokens):
   - **Incident Summary**: 2-3 sentence recap
   - **Key Findings by Perspective**: top 3 findings per analyst with ambiguity scores
   - **Recommendations**: all proposed recommendations
   - **FORCE PASS items**: analysts that didn't meet the threshold
2. Shut down completed analysts and DAs
3. Read `prompts/tribunal.md` for critic prompts
4. Replace placeholders:
   - `{FINDINGS_PACKAGE}` → compiled findings
   - `{TRIGGER_REASON}` → "User requested tribunal review"
   - `{INCIDENT_CONTEXT}` → Phase 0 details
5. Spawn UX Critic (Sonnet) + Engineering Critic (Opus) in parallel
6. Collect reviews, run consensus round:

| Level | Condition | Label |
|-------|-----------|-------|
| Strong | 2/2 APPROVE | `[Unanimous]` |
| Caveat | 1 APPROVE, 1 CONDITIONAL | `[Approved w/caveat]` |
| Split | 1+ REJECT | `[No consensus]` → user decision |

Split → share rationale, 1 final round only. Still split → present to user via `AskUserQuestion`.

7. Compile verdict, shut down critics, proceed to Phase 3.

### If Skip

Proceed directly to Phase 3.

→ **NEXT ACTION: Proceed to Phase 3 — Synthesis & Report.**

---

## Phase 3: Synthesis & Report

### Step 3.1

Integrate all verified analyst findings. Read from `.omc/state/incident-{short-id}/analyst-findings.md`.

### Step 3.2

Read `templates/report.md` and fill all sections with synthesized findings.

### Step 3.3

`AskUserQuestion`:
- "Is the analysis complete?"
- Options: "Complete" / "Need deeper investigation" / "Add recommendations" / "Share with team"

**Deeper investigation re-entry (max 2 loops):**

Before re-entry, increment `investigation_loops` counter in `.omc/state/incident-{short-id}/context.md`. If counter ≥ 2, inform user: "Maximum investigation depth reached. Proceeding with current findings." and auto-select "Complete".

1. Write current findings to `.omc/state/incident-{short-id}/analyst-findings.md`
2. Append iteration summary to `prior-iterations.md`
3. Identify gaps via `AskUserQuestion` (header: "Investigation Gaps"):
   - "Add new perspective" → spawn new analyst only (existing findings preserved)
   - "Re-examine with focus" → user specifies focus area → targeted follow-up tasks
4. New analyst runs → full Socratic DA + Scorer verification
5. Return to Phase 3 synthesis with expanded findings

→ **NEXT ACTION: Proceed to Phase 4 — Cleanup.**

---

## Phase 4: Cleanup

> Execute `../shared-v3/team-teardown.md`.
