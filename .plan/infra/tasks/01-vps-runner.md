# Task group — 01 VPS + self-hosted runner

**Depends on:** nothing (infrastructure bootstrap). **Blocks:** everything in this folder.

- [ ] **VPS provisioned** — a box with enough disk for the local-disk objects volume plus the root filesystem, MongoDB and Redis. Spec and OS per the plan (Linux; the runner labels assume x64).
- [ ] **Self-hosted GitHub Actions runner installed** — Settings → Actions → Runners → New self-hosted runner; the generated `./config.sh` run on the VPS (`infra/README.md` §Self-hosted runner).
- [ ] **Runner registered as a service** — `sudo ./svc.sh install && sudo ./svc.sh start` so it survives reboots.
- [ ] **Runner labeled** `self-hosted,linux,x64,bloberry` — so `release-desktop.yml`'s Linux job targets it specifically while macOS/Windows go to GitHub-hosted runners.
- [ ] **Sudoers rule, not blanket sudo** — the runner user gets exactly `sudo systemctl restart bloberry` via a sudoers entry (deploy needs it; nothing else does).
- [ ] **One runner serves everything** — confirmed as a monorepo (`architecture.md` §7): backend, web, CLI-Linux and desktop-Linux all deploy from this one registered runner.

**verification:** the runner shows online in repo settings; `sudo systemctl restart bloberry` works from the runner user; the runner survives a reboot.
