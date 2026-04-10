## Context

GitHub Pages is already enabled for `/docs` folder (used for `latest.txt`). The site will be available at `coleim.github.io/tsk`.

## Goals / Non-Goals

**Goals:**
- Simple, fast-loading landing page
- Clear installation instructions
- Screenshots showcasing the TUI
- Mobile-friendly design
- Easy to maintain (minimal dependencies)

**Non-Goals:**
- Complex build system
- Blog/news section
- User accounts or dynamic features
- JavaScript-heavy interactions

## Decisions

### 1. Plain HTML/CSS (no generator)
**Decision**: Use vanilla HTML/CSS with minimal JavaScript.
**Rationale**: Simple to maintain, no build step needed, instant deployment. The site is small (3-4 pages).
**Alternative considered**: Hugo/Jekyll - overkill for a small project site.

### 2. Single-page design with sections
**Decision**: One `index.html` with anchor links to sections (Features, Install, Screenshots).
**Rationale**: Simpler than multiple pages, better user experience for small content.

### 3. Dark theme matching TUI aesthetic
**Decision**: Dark background with accent colors matching the terminal app.
**Rationale**: Consistent branding, appeals to terminal users.

### 4. Screenshots as static images
**Decision**: PNG screenshots with optional animated GIF for demo.
**Rationale**: Simple, works everywhere, easy to update.

## Risks / Trade-offs

**[Trade-off]** No static site generator
→ Manual HTML editing, but content is minimal

**[Risk]** Screenshots may become outdated
→ Document how to capture screenshots in README or contributing guide
