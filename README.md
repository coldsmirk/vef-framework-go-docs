# VEF Framework Go Docs

Official documentation site for `vef-framework-go`, built with Docusaurus.

- Site: `https://coldsmirk.github.io/vef-framework-go-docs/`
- Default locale: Simplified Chinese
- English locale: `https://coldsmirk.github.io/vef-framework-go-docs/en/`

## Stack

- Docusaurus 3.9.2
- React 19
- TypeScript
- pnpm

## Local Development

Install dependencies:

```bash
pnpm install
```

Start the Chinese site locally:

```bash
pnpm start
```

Start the English site locally:

```bash
pnpm start:en
```

Start the Chinese locale explicitly:

```bash
pnpm start:zh
```

## Build And Preview

Build all locales:

```bash
pnpm build
```

Preview the production build:

```bash
pnpm serve
```

## Repository Structure

```text
docs/                                  # English source docs
i18n/zh-Hans/                          # Chinese translations
src/                                   # Docusaurus pages and theme code
static/                                # Static assets
docusaurus.config.ts                   # Site configuration
sidebars.ts                            # Sidebar config
```

## Notes

- The site is served at `/vef-framework-go-docs/`, with Chinese at the root path and English under `/en/`
- Use `pnpm build && pnpm serve` when checking multi-locale output
- This repository is the documentation site only; framework source code lives in the main `vef-framework-go` repository

## Deploy

Deploy with Docusaurus:

```bash
pnpm deploy
```
