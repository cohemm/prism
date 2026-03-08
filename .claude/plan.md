# Plan: Setup Agent → Team Member 전환 (incident-v2, plan-v2)

## 목표

현재 격리된 foreground subagent(setup-agent.md)로 실행되는 seed analysis를 **team member** 방식으로 전환:
- seed-analyst가 팀 멤버로 스폰되어 능동적 리서치 + seed analysis 수행
- 결과를 SendMessage로 team lead에게 전달
- incident-v2: Gather Evidence를 seed analysis에 통합 (코드베이스/MCP 능동 탐색)
- plan-v2: 동일 구조 전환 (리서치는 입력 파일/URL 분석 수준)

## 핵심 설계 결정

### AskUserQuestion 분리

team member는 AskUserQuestion을 사용할 수 없으므로:
- **Orchestrator가 처리**: 입력 수집, 심각도/상태, perspective 승인, ontology 소스 선택
- **Seed-analyst가 처리**: 능동적 리서치, 차원 평가, perspective 후보 생성

### Team 생성 시점 변경

- **현재**: Team은 Phase 1(analysts 스폰 시점)에서 생성
- **변경 후**: Team은 Phase 0.5(seed analysis 전)에서 생성 → seed-analyst도 team member로 참여

### Seed-analyst Agent Type

- `oh-my-claudecode:architect` (opus, READ-ONLY) — 리서치에 최적
- Worker Preamble 적용, SendMessage로 결과 보고
- 리서치 완료 후 shutdown

## 변경될 흐름

### incident-v2 새 흐름

```
Phase 0: Problem Intake (Orchestrator — NOT delegated)
  Step 0.1: Collect incident from user (or $ARGUMENTS)
  Step 0.2: Severity & Context (AskUserQuestion)
  Step 0.3: Fast Track Check
  Step 0.4: Generate short-id + create state dir

Phase 0.5: Team Creation + Seed Analysis
  Step 0.5.1: TeamCreate("incident-analysis-{short-id}")
  Step 0.5.2: Create seed-analyst task
  Step 0.5.3: Spawn seed-analyst (team member, background)
    → Active research: Grep codebase, Read files, git log, MCP (Sentry/Grafana)
    → Evaluate dimensions: domain, failure type, evidence, complexity, recurrence
    → Map to archetype candidates
    → Generate 3-6 perspective candidates + Perspective Quality Gate
    → SendMessage(perspectives + research summary) → team lead
  Step 0.5.4: Receive seed-analyst results
  Step 0.5.5: Shutdown seed-analyst
  [FAST_TRACK: Skip 0.5.2-0.5.5, lock 4 core perspectives directly]

Phase 0.6: Perspective Approval (Orchestrator)
  Step 0.6.1: AskUserQuestion for approval
  Step 0.6.2: Iterate until approved
  Step 0.6.3: Write perspectives.md
  [FAST_TRACK: Skip, perspectives already locked]

Phase 0.7: Ontology Scope Mapping (Orchestrator)
  → Execute shared/ontology-scope-mapping.md

Phase 0.8: Context & State Files
  Step 0.8.1: Write context.md (user input + seed research summary)
  Step 0.8.2: Write setup-complete.md sentinel (optional — for compatibility)
  Exit Gate: perspectives.md + context.md + ontology files written

Phase 1: Analyst Task Creation + Spawn (team already exists)
  Step 1.1: Create analyst tasks + DA task (with blockedBy)
  Step 1.2: Pre-assign owners
  Step 1.3: Spawn analysts in parallel (existing flow)
  Step 1.4: Spawn DA after analysts complete (existing flow)

Phase 2: Analysis Execution (unchanged)
Phase 2.5: Conditional Tribunal (unchanged)
Phase 3: Synthesis & Report (unchanged)
Phase 4: Cleanup (unchanged)
```

### plan-v2 새 흐름

```
Phase 0: Input Analysis (Orchestrator — NOT delegated)
  Step 0.1: Detect input type (file/URL/text/conversation)
  Step 0.2: Language detection
  Step 0.3: Extract planning context (goal/scope/constraints)
  Step 0.4: Fill gaps via AskUserQuestion
  Step 0.5: Generate short-id + create state dir

Phase 1: Team Creation + Seed Analysis
  Step 1.1: TeamCreate("plan-committee-{short-id}")
  Step 1.2: Create seed-analyst task
  Step 1.3: Spawn seed-analyst (team member, background)
    → Read input file/URL if provided
    → Evaluate dimensions: domain, complexity, risk, stakeholders, timeline, novelty
    → Generate 3-6 perspective candidates + Perspective Quality Gate
    → SendMessage(perspectives) → team lead
  Step 1.4: Receive seed-analyst results
  Step 1.5: Shutdown seed-analyst

Phase 1.5: Perspective Approval (Orchestrator)
  Step 1.5.1: AskUserQuestion for approval
  Step 1.5.2: Write perspectives.md

Phase 1.6: Ontology Scope Mapping (Orchestrator)
  → Execute shared/ontology-scope-mapping.md

Phase 1.7: Context & State Files
  Write context.md, setup-complete.md
  Exit Gate

Phase 2: Team Formation (tasks + owners — team already exists)
Phase 3: Parallel Analysis (unchanged)
Phase 4: DA Evaluation (unchanged)
Phase 5: Committee Debate (unchanged)
Phase 6: Plan Output (unchanged)
Phase 7: Cleanup (unchanged)
```

## 파일 변경 목록

### 수정할 파일

| # | 파일 | 변경 내용 |
|---|------|-----------|
| 1 | `skills/incident-v2/SKILL.md` | Phase 0~1 전면 재구성: 직접 intake + team-based seed analysis |
| 2 | `skills/incident-v2/docs/delegated-phases.md` | Phase 0에 능동적 리서치 단계 추가, setup-agent 위임 제거 |
| 3 | `skills/plan-v2/SKILL.md` | Phase 0~2 재구성: 직접 intake + team-based seed analysis |
| 4 | `skills/plan-v2/docs/delegated-phases.md` | setup-agent 위임 제거 |
| 5 | `skills/shared/setup-agent.md` | incident-v2/plan-v2 미사용 명시 (prd-v2만 사용) |

### 새로 생성할 파일

| # | 파일 | 내용 |
|---|------|------|
| 6 | `skills/incident-v2/prompts/seed-analyst.md` | Seed analyst 프롬프트 (능동적 리서치 + 차원 평가 + perspective 생성) |
| 7 | `skills/plan-v2/prompts/seed-analyst.md` | Seed analyst 프롬프트 (입력 분석 + 차원 평가 + perspective 생성) |

### 변경 없는 파일

| 파일 | 이유 |
|------|------|
| `skills/shared/da-evaluation-protocol.md` | 영향 없음 |
| `skills/shared/worker-preamble.md` | 그대로 seed-analyst에 적용 |
| `skills/shared/perspective-quality-gate.md` | 그대로 seed-analyst가 적용 |
| `skills/shared/ontology-scope-mapping.md` | Orchestrator가 직접 실행 (변경 없음) |
| `skills/prd-v2/*` | scope 밖 (기존 setup-agent.md 유지) |

## 구현 순서

1. **incident-v2 seed-analyst 프롬프트 작성** (`prompts/seed-analyst.md`)
   - Worker Preamble 호환 구조
   - 능동적 리서치 지시 (Grep, Read, Bash, MCP tools)
   - 차원 평가 + archetype 매핑 테이블
   - Perspective Quality Gate 적용
   - 출력: SendMessage로 perspectives + research summary

2. **incident-v2 SKILL.md 재구성**
   - Phase 0: Orchestrator 직접 intake (delegated-phases.md 참조 제거)
   - Phase 0.5: TeamCreate + seed-analyst 스폰
   - Phase 0.6: Perspective approval
   - Phase 0.7: Ontology mapping
   - Phase 0.8: Context + state files
   - Phase 1: Analyst tasks (team already exists)
   - Artifact persistence 테이블 업데이트
   - Gate Summary 업데이트

3. **incident-v2 delegated-phases.md 업데이트**
   - Phase 0에 능동적 리서치 포함
   - setup-agent 위임 문구 제거

4. **plan-v2 seed-analyst 프롬프트 작성** (`prompts/seed-analyst.md`)

5. **plan-v2 SKILL.md 재구성** (incident-v2와 동일 패턴)

6. **plan-v2 delegated-phases.md 업데이트**

7. **shared/setup-agent.md 업데이트** — incident-v2/plan-v2 미사용 명시

## Seed-analyst 프롬프트 핵심 구조 (incident-v2)

```markdown
# Seed Analyst — Incident Investigation

→ Apply worker preamble

## ROLE
You are the SEED ANALYST. Your job is to actively investigate the incident
and generate perspective candidates for the analysis team.

## INCIDENT CONTEXT
{INCIDENT_CONTEXT}

## SEVERITY: {SEVERITY}
## STATUS: {STATUS}
## EVIDENCE TYPES: {EVIDENCE_TYPES}

## PHASE 1: Active Research

MUST actively investigate using available tools:

| Evidence Type | Research Action |
|--------------|----------------|
| Error messages | Grep codebase for error strings |
| Stack traces | Read source files at referenced locations |
| Service names | Glob + Read for service configs and entry points |
| Recent deploys | Bash: git log --oneline --since="7 days ago" |
| Metrics/dashboards | ToolSearch for Grafana/Sentry MCP → query |
| Logs | ToolSearch for log MCP (Loki, ClickHouse) → query |

Write detailed research findings as you go.

## PHASE 2: Seed Analysis

Evaluate across 5 dimensions using research findings:
[dimension table from delegated-phases.md]

Map to archetype candidates using:
[characteristic-to-archetype mapping table]

## PHASE 3: Generate Perspectives

Generate 3-6 orthogonal perspectives.
Apply Perspective Quality Gate (shared/perspective-quality-gate.md).

## OUTPUT (via SendMessage to team-lead)

### Research Summary
[Key findings from active investigation — evidence discovered, files examined]

### Perspectives
[Per perspective: ID, Name, Scope, Key Questions, Model, Agent Type, Rationale]
```

## 리스크 & 완화

| 리스크 | 완화 |
|--------|------|
| Seed-analyst가 너무 오래 리서치 | 프롬프트에 시간 제한 지시: "Max 5 minutes active research. Prioritize high-signal evidence." |
| MCP 도구 없는 환경 | "If no MCP tools available, skip MCP queries. Investigate using codebase tools only." |
| SendMessage payload 너무 큼 | perspectives만 SendMessage, 상세 research는 state file에 Write |
| Fast Track에서 seed-analyst 불필요 | FAST_TRACK은 seed-analyst 스킵, 4 core 직접 잠금 (기존과 동일) |
