# Incident RCA Report

---

## Executive Summary

[2-3 sentences: what happened, root cause, severity, current status]

## Incident Overview

- **Incident**: {incident title}
- **Analysis Date**: {date}
- **Method**: {N}-perspective multi-perspective analysis + Socratic verification
- **Reference Docs**: {number of ontology documents} documents referenced
- **Perspectives Used**: {list of perspectives}

---

## Timeline

| Time | Event | Source |
|------|-------|--------|
| {time} | {event description} | {how this was determined} |

> If timeline data is insufficient, state: "Timeline could not be fully reconstructed from available evidence."

> **Timeline Reconstruction Rules** (the synthesizer MUST follow these when populating the table above):
>
> 1. **Chronological ordering**: Events MUST be listed in strict chronological order. If exact timestamps are unavailable, use relative ordering (e.g., "Before deploy", "T+0", "T+5min") and note the precision level.
>
> 2. **Required milestone events**: Actively search findings for these milestone events and include them when evidence exists:
>    - **Trigger point**: The immediate event that initiated the incident (e.g., deployment, config change, traffic spike)
>    - **First impact**: When users or systems first experienced the effect
>    - **Detection**: When the incident was detected (alert fired, user report, manual discovery)
>    - **Escalation points**: Key moments where impact worsened or scope expanded
>    - **Mitigation start**: When remediation actions began
>    - **Resolution**: When the incident was resolved or mitigated
>
> 3. **Source attribution**: Every timeline entry MUST cite how the event was determined — e.g., "deployment log", "alert timestamp", "user report", "inferred from code change history", or "{perspective name} analysis". Never list an event without a source.
>
> 4. **Cross-perspective correlation**: When multiple perspectives reference temporal events, cross-reference them to build the most complete timeline. If two perspectives disagree on timing, include both with a note (e.g., "~14:00 per infra analysis, ~14:05 per UX reports").
>
> 5. **Do not fabricate**: If temporal data is sparse, include only events with evidence. A timeline with 2-3 well-sourced entries is better than a speculative 10-entry timeline. Mark uncertain entries with `(approximate)` or `(inferred)`.

---

## Root Cause

### Primary Root Cause

**What**: {description of the root cause}

**Where**: {code path / system component}

**Why**: {explanation of why this caused the incident}

**Confidence**: {Verified / Partial / Unverified}

### Contributing Factors

| # | Factor | Relationship to Root Cause | Confidence |
|---|--------|---------------------------|------------|
| 1 | {factor} | {how it contributed} | {badge} |

### Trigger

{The immediate event that initiated the incident}

---

## UX Impact Analysis

### User Experience During Incident

{What users saw, experienced, or were unable to do}

### Affected User Flows

| # | Flow | Impact | Affected Users |
|---|------|--------|---------------|
| 1 | {user flow} | {what broke / degraded} | {scope: all users / segment / %} |

### Technical Cause → UX Effect Mapping

| Technical Cause | Data Flow Path | UX Effect | User Saw | Severity |
|----------------|---------------|-----------|----------|----------|
| {code/system issue} | {component A → component B → … → UI layer} | {what user experienced} | {exact symptom: error message, blank screen, stale data, etc.} | {CRITICAL/HIGH/MEDIUM} |

> **Cause-Effect Mapping Rules** (the synthesizer MUST follow these when populating the table above):
>
> 1. **Trace backward from symptom**: Start from what the user saw (error message, broken UI, missing data) and trace backward through the stack: UI → API → service → database/infra. Every row MUST include the full data flow path, not just the endpoints.
>
> 2. **Never stop at the frontend**: If the UX effect originates from a backend or infrastructure failure, the `Technical Cause` column MUST reference the deepest root (e.g., database timeout, misconfigured deployment), NOT the frontend symptom (e.g., "React component crashed"). The frontend is the **effect**, not the cause.
>
> 3. **Distinguish direct vs. cascading effects**: If a single technical cause produced multiple UX effects through different data flow paths, list each as a separate row. If multiple technical causes converged to produce one UX effect, list each cause as a separate row pointing to the same UX effect.
>
> 4. **Include timing and conditions**: When available, annotate each mapping with the conditions under which the symptom appears (e.g., "only for users with > 100 items", "only during peak hours", "after retry exhaustion at 30s").
>
> 5. **Severity must reflect user impact, not technical severity**: A critical backend error that is gracefully handled by the frontend is MEDIUM. A minor API change that breaks a core user flow is CRITICAL. Always assess severity from the user's perspective.

### UX Effect Data Flow Traces

> For the top 1-3 most severe cause-effect mappings above, provide a detailed data flow trace showing exactly how the technical failure propagated to the user. Use the format below:

```
[Symptom] User saw: {exact symptom}
  ← [UI Layer] {what the frontend did/failed to do}
    ← [API Layer] {API response or failure}
      ← [Service Layer] {service behavior}
        ← [Root] {infrastructure/code/data issue}
```

> **Trace rules**:
> - Each trace MUST start from the user-visible symptom and end at the deepest identified root cause.
> - If a layer is not involved, skip it — but never skip a layer that IS involved just because evidence is weak. Mark uncertain layers with `(inferred)`.
> - If the trace cannot be fully reconstructed, end with `← [Unknown] Evidence insufficient beyond this point` and note the gap.

---

## Per-Perspective Analysis Summary

### {Perspective Name}
- **Scope**: {what this perspective examined}
- **Findings**: {N} items (CRITICAL: {n}, HIGH: {n}, MEDIUM: {n})
- **Key Finding**: {1-2 sentence summary}
- **Verification Score**: {score} ({verdict})

---

## Cross-Perspective Analysis

> **Synthesis Rules**: The following subsections MUST be populated by correlating findings across ALL perspectives. Do not simply list each perspective's conclusions independently — actively compare, cross-reference, and integrate them.

### Corroborated Root Causes

Items where **2+ perspectives independently identified the same cause**. These carry the highest confidence.

| Root Cause | Agreeing Perspectives | Combined Confidence | Evidence Summary |
|------------|----------------------|--------------------|--------------------|
| {root cause} | {perspective A, perspective B, ...} | {highest badge among agreeing perspectives} | {how each perspective arrived at the same conclusion independently} |

> **Correlation rule**: When multiple perspectives identify the same root cause using different evidence paths (e.g., code-level analysis finds a bug AND infrastructure analysis finds the same failure mode), mark as **Strongly Corroborated**. When they share the same evidence chain, mark as **Corroborated**.

### Conflicting Conclusions

Items where perspectives reached **different or contradictory root cause conclusions**. Each conflict MUST include an explicit resolution.

| Finding | Perspective A View | Perspective B View | Resolution |
|---------|-------------------|-------------------|------------|
| {finding} | {perspective A's conclusion} | {perspective B's conclusion} | {which is more likely correct and why, or how both are partially correct} |

> **Resolution rule**: Prefer the conclusion from the perspective with the higher Socratic verification score. If scores are equal, prefer the perspective whose scope is most directly relevant to the finding. If unresolvable, state both conclusions and flag for manual investigation.

### Integrated Insights

Issues only visible when **combining findings across perspectives** — no single perspective could have identified these alone.

| Insight | Contributing Perspectives | How Perspectives Combine |
|---------|--------------------------|-------------------------|
| {insight} | {perspectives involved} | {what each perspective contributed and why the combined view reveals something new} |

> **Integration rule**: Actively look for these patterns:
> - **Technical cause → UX effect chains**: A code/infrastructure perspective identifies the technical failure, while a UX perspective reveals how it manifested to users. Combine into a single cause-effect narrative.
> - **Blast radius amplification**: One perspective identifies the bug, another reveals that the blast radius was larger than the bug alone would suggest (e.g., cascading failures, retry storms).
> - **Hidden dependencies**: One perspective identifies a component failure, another reveals an undocumented dependency that propagated the failure to unrelated features.
> - **Temporal correlation**: One perspective identifies what changed, another identifies when the impact began — together they confirm or refute causality.

---

## Action Items

> **Action Items Generation Rules** (the synthesizer MUST follow these when populating the subsections below):
>
> 1. **Extract from findings, not assumptions**: Every action item MUST trace back to a specific finding from the analysis. Do not add generic best-practice items unless they directly address an identified gap.
>
> 2. **Prioritize by impact × effort**: Within each category, order items by estimated impact (how much it reduces recurrence risk) balanced against implementation effort. Highest impact, lowest effort items come first.
>
> 3. **Be specific and actionable**: Each item MUST include enough context to act on — reference specific code paths, system components, configuration files, or monitoring gaps. Avoid vague items like "improve monitoring" — instead: "Add latency P99 alert on `/api/checkout` endpoint with threshold > 2s".
>
> 4. **Category assignment rules**:
>    - **Immediate fixes**: Direct code/config changes that prevent this exact incident from recurring. These should be deployable within 1-2 days.
>    - **Short-term improvements**: Systemic improvements (monitoring, alerting, testing, runbooks) that improve detection or reduce blast radius. Target: 1-2 weeks.
>    - **Long-term considerations**: Architectural changes, process improvements, or design patterns that address the structural weakness. Target: 1+ months.
>
> 5. **Cross-reference with root cause**: At least one immediate fix MUST directly address the primary root cause. If the root cause cannot be immediately fixed, explain why and ensure a short-term mitigation is listed instead.
>
> 6. **Include ownership hints**: When possible, indicate which team or system area each action item belongs to (e.g., "[Backend]", "[Platform]", "[Frontend]", "[SRE]").

### Immediate Fixes (Prevent Recurrence)
- [ ] {specific action with code/system reference}

### Short-Term Improvements (Monitoring / Alerting / Testing)
- [ ] {specific improvement}

### Long-Term Considerations (Architecture / Process)
- [ ] {specific consideration}

---

## Confidence Summary (Socratic Verification)

| Perspective | Verification Rounds | Score | Verdict |
|-------------|-------------------|-------|---------|
| {name} | {rounds} | {score} | {verdict} |

---

## Appendix

### Referenced Documents
| # | Document Path | Related Perspective | Reference Count |
|---|--------------|-------------------|-----------------|

### Raw Evidence
{Key code snippets, log entries, or data points that support the root cause determination}

---

## Badge Mapping Rules

> These rules govern how badges are rendered throughout this report.
> The synthesizer MUST apply these mappings consistently to all tables and inline references.

### Confidence Badges (Socratic Verification Score → Badge)

| Condition | Badge | Rendering |
|-----------|-------|-----------|
| score >= 0.7 | **Verified** | ✅ Verified |
| 0.4 <= score < 0.7 | **Partial** | ⚠️ Partial |
| score < 0.4 | **Unverified** | ❓ Unverified |
| score unavailable | **Unverified** | ❓ Unverified |

- Each finding inherits the confidence badge of its parent perspective.
- If Socratic verification data was unavailable for all perspectives, add this note under the Confidence Summary table:
  > "Socratic verification data was unavailable. All confidence levels are marked as 'Unverified'."

### Severity Badges (Impact Level → Badge)

| Level | Badge | Criteria |
|-------|-------|----------|
| CRITICAL | 🔴 CRITICAL | Service outage, data loss, or security breach affecting all/most users |
| HIGH | 🟠 HIGH | Major feature broken or significant degradation for a user segment |
| MEDIUM | 🟡 MEDIUM | Partial degradation, workaround available, or limited user impact |
| LOW | 🟢 LOW | Minor issue, cosmetic, or negligible user impact |

### Badge Application Points

- **Contributing Factors table** (`## Root Cause > ### Contributing Factors`): Use confidence badges in the `Confidence` column.
- **Technical Cause → UX Effect Mapping table** (`## UX Impact Analysis`): Use severity badges in the `Severity` column.
- **Per-Perspective Analysis Summary** (`## Per-Perspective Analysis Summary`): Use confidence badges for the `Verification Score` field, rendering as `{score} ({badge})`.
- **Confidence Summary table** (`## Confidence Summary`): Use confidence badges in the `Verdict` column.
- **Findings counts** in Per-Perspective summaries: Prefix counts with severity badges (e.g., `🔴 CRITICAL: 2, 🟠 HIGH: 3, 🟡 MEDIUM: 1`).
