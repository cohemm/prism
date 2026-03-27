# Ontology Scope Schema & Text Block

Read this file only when executing Step 4 (writing `ontology-scope.json`) or Phase B (generating `{ONTOLOGY_SCOPE}` text block).

## Step 4: Write `ontology-scope.json`

Combine all collected sources into `{STATE_DIR}/ontology-scope.json`.

If no sources and `{AVAILABILITY_MODE}`=`optional` → warn and set `{ONTOLOGY_SCOPE}` = "N/A". If `required` → error and **STOP**.

### Schema

```json
{
  "sources": [
    {
      "id": 1,
      "type": "doc",
      "path": "/path/to/docs",
      "domain": "inferred domain",
      "summary": "Documentation directory",
      "key_topics": ["topic1", "topic2"],
      "status": "available",
      "access": {
        "tools": ["Read"],
        "instructions": "Use the Read tool with offset/limit to read files in the directory."
      }
    },
    {
      "id": 2,
      "type": "mcp_query",
      "server_name": "grafana",
      "domain": "monitoring",
      "summary": "Grafana monitoring dashboards and metrics",
      "key_topics": ["prometheus", "loki", "dashboards"],
      "status": "available",
      "access": {
        "tools": ["mcp__grafana__query_prometheus", "mcp__grafana__query_loki_logs"],
        "instructions": "Call ToolSearch(query=\"select:mcp__{server_name}__{tool_name}\") to load each tool before use, then call directly.",
        "capabilities": "Query Prometheus metrics, Loki logs, dashboards",
        "getting_started": "Start with list_datasources to discover available data",
        "error_handling": "If a tool call fails, note the error and continue. Do NOT retry more than once.",
        "safety": "SELECT/read-only queries only"
      }
    },
    {
      "id": 3,
      "type": "web",
      "url": "https://example.com/docs",
      "domain": "API documentation",
      "summary": "1-2 line summary",
      "key_topics": ["keyword1", "keyword2", "keyword3"],
      "status": "available",
      "access": {
        "instructions": "Content summary provided. Use WebFetch for deeper exploration.",
        "cached_summary": "Fetched content summary from pool build time"
      }
    },
    {
      "id": 4,
      "type": "file",
      "path": "/path/to/file",
      "domain": "domain",
      "summary": "1-2 line summary",
      "key_topics": ["keyword1", "keyword2"],
      "status": "available",
      "access": {
        "instructions": "Content summary provided. Original file at path.",
        "cached_summary": "Read content summary from pool build time"
      }
    },
    {
      "id": 5,
      "type": "web",
      "url": "https://failed.example.com",
      "status": "unavailable",
      "reason": "fetch failed: 404"
    }
  ],
  "totals": {
    "doc": 1,
    "mcp_query": 1,
    "web": 1,
    "file": 1,
    "unavailable": 1
  },
  "citation_format": {
    "doc": "source:section",
    "web": "url:section",
    "file": "file:path:section",
    "mcp_query": "mcp-query:server:detail"
  }
}
```

### Field Rules

- `sources[].id`: sequential integer, starts at 1
- `sources[].type`: one of `doc`, `mcp_query`, `web`, `file`
- `sources[].status`: `available` or `unavailable`
- `sources[].key_topics`: 3-5 keywords (inferred or extracted at pool build time)
- `sources[].access`: present only when `status == "available"`
- `sources[].reason`: present only when `status == "unavailable"`
- `totals`: count per type. `unavailable` counts all failed sources regardless of type

---

## Phase B: Generate `{ONTOLOGY_SCOPE}` Text Block

The orchestrator reads `{STATE_DIR}/ontology-scope.json` and generates a text block for analyst prompt injection. Only `available` sources are included.

### Text block format

```
Your reference documents and data sources:

- doc: {summary} ({status})
  Directories: {path}
  Access: {access.instructions}
    {access.tools — one per line}

- mcp-query: {server_name}: {summary}
  Tools (read-only): {access.tools}
  Access: {access.instructions}
  Capabilities: {access.capabilities}
  Getting started: {access.getting_started}
  Error handling: {access.error_handling}

- web: {url}: {domain} — {summary}
  Access: {access.instructions}
  {access.cached_summary}

- file: {path}: {domain} — {summary}
  Access: {access.instructions}
  {access.cached_summary}

Explore these sources through your perspective's lens.
Cite findings as: {citation_format values}.
```

**Backward compatibility:** If `ontology-scope.json` does not exist at read time, inject: `{ONTOLOGY_SCOPE}` = "N/A — ontology scope file not found. Analyze using available evidence only."
