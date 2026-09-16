# Loop Runs

This directory contains the persistent history of Loop Engineering executions.

Each run should have its own directory:

```text
runs/
└── <run-id>/
    ├── goal.md
    ├── state.md
    ├── decisions.md
    ├── failures.md
    └── result.md
```

Run history is append-oriented engineering evidence.

A completed run should remain available for future analysis, debugging, and improvement of the Loop Engineering system.

The current active run is represented by:

```text
../state.md
```

Historical runs must not be treated as the current project state.
