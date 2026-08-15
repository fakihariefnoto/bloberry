# Screen — search

## Purpose & context

- **User goal**: find a file by name, type, size or date within the current tenant (PRD S6). Mobile's quick lookup — "where did I put that photo".
- **Entry points**: search icon on `files` app bar (`search` route, `q: String?` arg pre-fillable).
- **Exit points**: tap a result → `file-detail`; back → `files`. Results are tenant-wide; permission filtering applies server-side (a viewer sees only what they can read).
- **Data needed**: debounced query → `GET` search (name, type, size range, date); recent searches (local), suggestions (recently accessed files the user can read).

## States

- [x] Pre-query — recent searches + suggestions (never a blank waiting screen)
- [x] Loading (debounced results, subtle inline indicator)
- [x] Populated results
- [x] Empty — "no results for <query>"
- [x] Error (inline, retry)
- [x] Domain-specific — query too short (min length hint)

## Style reference

- **Components used**: search screen pattern — top search field (autofocus), recent/suggested lists pre-query, results list post-query. Row = type icon + name + parent folder path + size/date. No tab bar (pushed over `files`), back in the leading position.
- No token deltas.

## Wireframe — mobile (pre-query)

```
┌─────────────────────────────┐
│  ←  [ 🔍 search files... ]  │
├─────────────────────────────┤
│  RECENT                     │
│  ┌────────────────────────┐ │
│  │ 🕐 q3-report.pdf       │ │
│  └────────────────────────┘ │
│  ┌────────────────────────┐ │
│  │ 🕐 hero.png            │ │
│  └────────────────────────┘ │
│                             │
│  SUGGESTED                  │
│  ┌────────────────────────┐ │
│  │ 📄 index.html  · Web/  │ │
│  └────────────────────────┘ │
│  ┌────────────────────────┐ │
│  │ 🖼 hero.png     · 2026/ │ │
│  └────────────────────────┘ │
│  ┌────────────────────────┐ │
│  │ 📁 _assets     · Root/ │ │
│  └────────────────────────┘ │
└─────────────────────────────┘
```

## Wireframe — mobile (results)

```
┌─────────────────────────────┐
│  ←  [ hero ]                │
├─────────────────────────────┤
│  2 results                  │
│  ┌────────────────────────┐ │
│  │ 🖼 hero.png     2.4 MB │  │
│  │   Web/2026 · public   │  │
│  └────────────────────────┘ │
│  ┌────────────────────────┐ │
│  │ 🖼 hero-2x.png  4.8 MB │  │
│  │   Web/2026 · private  │  │
│  └────────────────────────┘ │
└─────────────────────────────┘
```

## Interactions

- **Pre-query** shows recent searches (tap to re-run) and suggested files the user can read — the "something useful before typing" rule (`design-collection/mobile-screen/patterns.md` Search).
- **Debounce** ~300ms; results replace recent/suggested content entirely (no stale mix). Loading is a subtle inline indicator, not a blocking spinner.
- **Result row** shows the parent folder path (the result's `ancestors` breadcrumb) — a search hit without context is "found it but where is it".
- **Empty**: "No results for 'hero' · Try a different name, type or date" — a real copy, never "No data".
- **Query too short** (<2 chars): hint caption, no request fired.
- **Back** → `files`, clearing the query. Search state is not persisted across the tab bar.
- **A11y**: the field autofocuses with the keyboard up; results announce count ("2 results"); each row is one tappable target ≥48dp.
