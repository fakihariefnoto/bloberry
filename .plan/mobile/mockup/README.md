# Mockups — Mobile

Each screen is a `.md` file. This is a design-thinking artifact, not a decoration — a wireframe full of generic "Heading Text" placeholders hasn't actually been designed yet. Every screen file has:

1. **Purpose & context** — the user's goal, entry/exit points, and the actual data this screen needs.
2. **States** — every state the screen can be in (loading/empty/populated/error/domain-specific), not just the happy path.
3. **Style reference** — deltas/exceptions from `../../design/style-guide.md` only (the shared doc, also used by web).
4. **Wireframe(s)** — one per state that's meaningfully different in layout, ASCII art in a fenced code block, using this screen's real content.
5. **Interactions** — what each element does, validation rules, loading/disabled behavior.

See `screen-example.md` for the format, and use the `refine-mockup` skill to iterate on an existing screen.
