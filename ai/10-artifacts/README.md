# Engineering Artifacts

## Purpose

This directory contains generated engineering artifacts produced during the development lifecycle.

Artifacts are temporary or evolving documents created to support engineering activities.

They are not considered architectural truth unless explicitly promoted into:

- architecture documentation
- standards
- bounded context documentation
- Architecture Decision Records

---

# Artifact Philosophy

The distinction is:
Knowledge

=
Permanent understanding

Artifacts

=
Work produced while creating knowledge


Artifacts help engineers think.

Documentation preserves final decisions.

---

# Artifact Categories

## Proposals

Location:


proposals/


Contains:

- feature proposals
- architecture proposals
- implementation plans
- migration proposals

Examples:

- new bounded context proposal
- API redesign proposal
- storage strategy proposal

---

## Research

Location:


research/


Contains technical investigations.

Examples:

- technology comparisons
- proof of concept results
- protocol analysis
- performance experiments

Research should answer:

- What problem are we solving?
- What options exist?
- What are the trade-offs?

---

## Reviews

Location:


reviews/


Contains engineering evaluations.

Examples:

- architecture reviews
- security reviews
- code reviews
- design reviews

Reviews identify:

- risks
- improvements
- approval status

---

## Reports

Location:


reports/


Contains analytical documents.

Examples:

- project reports
- technical assessments
- progress reports
- incident reports

---

## Release Notes

Location:


release-notes/


Contains user-facing and internal release documentation.

Examples:

- version summaries
- migration instructions
- breaking changes
- upgrade notes

---

# Lifecycle

Artifacts follow this lifecycle:


Created

↓

Reviewed

↓

Validated

↓

Promoted or Archived


---

# Promotion Rules

An artifact becomes permanent knowledge only when necessary.

Examples:

A proposal becomes:


07-decisions/adr-xxxx.md


if it creates an architectural decision.

A research document becomes:


03-standards


if it defines an engineering standard.

A review becomes:


08-agent-memory


if it represents an active constraint.

---

# AI Rules

AI assistants should:

- create artifacts when useful
- keep temporary reasoning separate from permanent knowledge
- reference artifacts during discussions
- promote important conclusions

AI assistants should not:

- modify permanent architecture documentation without justification
- store temporary thoughts as project truth
- create unnecessary documents

---

# Success Criteria

A healthy artifact system provides:

- traceability
- transparency
- historical context
- better decision making

Artifacts are the working memory of engineering.