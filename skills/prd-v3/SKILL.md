---
name: prd-v3
description: PRD 정책 충돌 분석 — prd.md 파일을 입력받아 코드베이스 정책 문서 대비 충돌/모호성을 다관점 분석하고 PM용 리포트를 생성합니다. "prd 분석", "PRD 정책 검토", "기획서 리뷰", "prd 충돌 확인", "PRD policy review" 등의 요청에 사용하세요. 기획서나 PRD를 언급하며 정책/충돌/모호성/리뷰를 요청하면 반드시 이 스킬을 사용합니다.
version: 1.0.0
user-invocable: true
allowed-tools: Skill, Task, Read, Write, Bash, Glob, Grep, AskUserQuestion, ToolSearch
---

# PRD Policy Analysis (Wrapper for analyze)

PRD 파일을 입력받아 코드베이스의 정책 문서(ontology)와 대조하여 정책 충돌과 모호성을 찾아내는 스킬. 내부적으로 `prism:analyze` 스킬에 분석을 완전 위임하고, 그 결과물을 PM이 읽기 쉬운 형태로 후처리합니다.

## Prerequisite

> Read and execute `../shared-v3/prerequisite-gate.md`. Set `{PROCEED_TO}` = "Phase 0".

---

## Phase 0: Input

### Step 0.1: Get PRD File Path

`$ARGUMENTS`에서 PRD 파일 경로를 추출합니다.

- 경로가 있으면 → `Read`로 파일 존재 확인
- 경로가 없으면 → `AskUserQuestion` (header: "PRD 파일", question: "분석할 PRD 파일 경로를 알려주세요.")
- 파일이 없으면 → ERROR: "PRD 파일을 찾을 수 없습니다: {path}"

### Step 0.2: Generate Session ID

```bash
uuidgen | tr '[:upper:]' '[:lower:]' | cut -c1-8
```

Generate ONCE, reuse throughout. Create state directory:

```bash
mkdir -p ~/.prism/state/prd-{short-id}
```

### Step 0.3: Language Detection

1. CLAUDE.md에 `Language` 지시가 있으면 → 해당 언어
2. 없으면 → 유저 입력 언어 감지
3. `{REPORT_LANGUAGE}`로 저장

### Phase 0 Exit Gate

- [ ] PRD 파일 경로 확인 및 파일 존재 검증
- [ ] `{short-id}` 생성 및 `~/.prism/state/prd-{short-id}/` 디렉토리 생성
- [ ] `{REPORT_LANGUAGE}` 결정

→ **NEXT: Phase 1 — Config 생성 및 analyze 호출**

---

## Phase 1: Config & Analyze 호출

### Step 1.1: Read PRD

PRD 파일 전문을 `Read`로 읽습니다. PRD의 제목, 주요 기능 요구사항(FR), 비기능 요구사항(NFR)을 파악합니다.

같은 디렉토리에 관련 파일(handoff, constraints 등)이 있으면 함께 읽습니다.

### Step 1.2: Create Analyze Config

다음 JSON을 `~/.prism/state/prd-{short-id}/analyze-config.json`에 작성합니다:

```json
{
  "topic": "PRD 정책 충돌 분석: {PRD 제목} — 이 PRD가 코드베이스의 기존 정책에 위배되거나 모호한 부분이 있는지 다관점으로 분석",
  "input_context": "{PRD 파일 절대경로}",
  "report_template": "{이 스킬의 절대경로}/templates/report.md",
  "seed_hints": "먼저 {PRD 파일 절대경로}에 있는 PRD 파일을 Read하라. PM(기획자) 관점에서 정책 도메인을 추출하라. 엔지니어링 구현 세부사항이 아니라 비즈니스 정책 충돌, 규칙 모순, 정의되지 않은 엣지케이스, 모호한 요구사항에 집중하라. PRD의 각 기능 요구사항이 기존 정책 문서와 충돌하는지, 기존 정책에서 다루지 않는 새로운 영역인지 분류하라.",
  "ontology_mode": "required"
}
```

> `{이 스킬의 절대경로}`는 이 SKILL.md가 위치한 디렉토리의 절대경로입니다. `Bash`로 확인: 이 SKILL.md 파일 경로에서 디렉토리를 추출합니다.

### Step 1.3: Invoke Analyze

```
Skill(skill="prism:analyze", args="--config ~/.prism/state/prd-{short-id}/analyze-config.json")
```

analyze가 완료될 때까지 대기합니다. analyze는 내부적으로:
- seed analyst로 PRD와 정책 도메인 조사
- 다관점 perspective 생성 및 유저 승인
- 각 perspective별 analyst 스폰 (정책 충돌 분석)
- Socratic verification으로 발견사항 검증
- 리포트 생성

### Step 1.3.1: Snapshot Before Analyze

analyze 호출 **직전에** 기존 analyze 디렉토리 목록을 스냅샷합니다:

```bash
ls -d ~/.prism/state/analyze-* 2>/dev/null > ~/.prism/state/prd-{short-id}/analyze-dirs-before.txt || touch ~/.prism/state/prd-{short-id}/analyze-dirs-before.txt
```

### Step 1.4: Locate Analyze Output

analyze 완료 후, 호출 전후 디렉토리를 비교하여 새로 생긴 analyze 디렉토리를 찾습니다:

```bash
comm -13 <(sort ~/.prism/state/prd-{short-id}/analyze-dirs-before.txt) <(ls -d ~/.prism/state/analyze-* 2>/dev/null | sort)
```

새로 생긴 디렉토리가 정확히 1개여야 합니다. 0개면 ERROR: "analyze가 state 디렉토리를 생성하지 않았습니다." 2개 이상이면 가장 최근 것을 선택합니다.

해당 디렉토리에서 다음 파일들을 확인합니다:
- `analyst-findings.md` — 검증된 분석 결과
- `verification-log.json` — Socratic 검증 점수 (없을 수 있음)

이 경로를 `{ANALYZE_STATE_DIR}`로 저장합니다.

### Phase 1 Exit Gate

- [ ] `analyze-config.json` 작성 완료
- [ ] `prism:analyze` 스킬 호출 완료
- [ ] `{ANALYZE_STATE_DIR}` 식별 및 `analyst-findings.md` 존재 확인

→ **NEXT: Phase 2 — 후처리 (PM용 리포트 변환)**

---

## Phase 2: Post-Processing (PM용 리포트 생성)

analyze가 생성한 결과물은 기술적 분석 리포트입니다. PM이 바로 활용할 수 있도록 후처리 에이전트가 변환합니다.

### Step 2.1: Spawn Post-Processor Agent

Read `prompts/post-processor.md` (relative to this SKILL.md).

```
Task(
  subagent_type="oh-my-claudecode:analyst",
  model="opus",
  prompt="{post-processor prompt with placeholders replaced}"
)
```

**CRITICAL: Do NOT add `run_in_background=true`.** 후처리 결과를 기다려야 합니다.

Placeholder replacements:
- `{ANALYZE_STATE_DIR}` → Phase 1.4에서 식별한 analyze 결과 디렉토리 경로
- `{PRD_FILE_PATH}` → PRD 파일 절대경로
- `{PRD_STATE_DIR}` → `~/.prism/state/prd-{short-id}`
- `{REPORT_LANGUAGE}` → Phase 0.3에서 결정한 언어
- `{SHORT_ID}` → 세션 ID

### Step 2.2: Verify Output

후처리 에이전트 완료 후, 리포트 파일 존재를 확인합니다:

```
~/.prism/state/prd-{short-id}/prd-policy-review-report.md
```

파일이 없으면 ERROR: "후처리 에이전트가 리포트를 생성하지 못했습니다."

### Phase 2 Exit Gate

- [ ] 후처리 에이전트 완료
- [ ] `prd-policy-review-report.md` 존재
- [ ] 리포트에 "PM 의사결정 체크리스트" 섹션 존재 확인 (`Grep`)

→ **NEXT: Phase 3 — 리포트 전달**

---

## Phase 3: Output

### Step 3.1: Copy Report to PRD Directory

PRD 파일이 있는 디렉토리에도 리포트 사본을 저장합니다:

```bash
cp ~/.prism/state/prd-{short-id}/prd-policy-review-report.md {PRD_DIR}/prd-policy-review-report.md
```

### Step 3.2: Report to User

유저에게 결과를 알립니다:

```
PRD 정책 분석이 완료되었습니다.

리포트 위치:
- {PRD_DIR}/prd-policy-review-report.md
- ~/.prism/state/prd-{short-id}/prd-policy-review-report.md

analyze 원본 결과: {ANALYZE_STATE_DIR}/
```
