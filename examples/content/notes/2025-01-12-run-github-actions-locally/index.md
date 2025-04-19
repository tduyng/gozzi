+++
title = "[Note] Run Github actions locally"
date = 2025-01-12
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Run Github actions locally

Run GitHub actions locally with [act](https://github.com/nektos/act) to skip the long feedback loop! It uses Docker to pull images and run workflows right on your machine.

```bash
# Run the `push` event
act

# Run a specific event or job
act pull_request
act -j test_unit
```

Need a GitHub token? Use the `-s` flag with the GitHub CLI:

```bash
act -s GITHUB_TOKEN="$(gh auth token)"
```

Quick, easy, and no more waiting for GitHub to run workflows!
