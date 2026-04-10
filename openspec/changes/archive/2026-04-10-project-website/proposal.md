## Why

The project needs a dedicated website to showcase features, provide installation instructions, and display screenshots. Currently, the only documentation is the README on GitHub. A static site hosted on GitHub Pages will improve discoverability and user onboarding.

## What Changes

- Create static website using a simple static site generator (Hugo or plain HTML/CSS)
- Add landing page with feature highlights
- Add installation guide page
- Add screenshots/demo section
- Configure GitHub Pages deployment from `docs/` folder
- Update `docs/latest.txt` workflow to not conflict with site

## Capabilities

### New Capabilities
- `static-website`: Static site structure, pages, and assets in docs/ folder

### Modified Capabilities
<!-- No spec-level changes to existing behavior -->

## Impact

- **docs/**: Will contain website files (currently only has `latest.txt`)
- **GitHub Pages**: Already configured for /docs folder - site will be at coleim.github.io/tsk
- **CI**: May need workflow to build site if using generator
