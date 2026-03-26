import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";
import { themes as prismThemes } from "prism-react-renderer";

const config: Config = {
  title: "VEF Framework Go",
  tagline:
    "Resource-driven Go APIs with FX, Fiber, and built-in enterprise capabilities",
  favicon: "img/favicon.svg",

  future: {
    v4: true,
  },

  url: "https://coldsmirk.github.io",
  baseUrl: "/vef-framework-go-docs/",

  organizationName: "coldsmirk",
  projectName: "vef-framework-go-docs",

  onBrokenLinks: "throw",

  markdown: {
    mermaid: true,
    hooks: {
      onBrokenMarkdownLinks: "warn",
    },
  },

  themes: ["@docusaurus/theme-mermaid"],

  i18n: {
    defaultLocale: "en",
    locales: ["en", "zh-Hans"],
    localeConfigs: {
      en: {
        label: "English",
        direction: "ltr",
      },
      "zh-Hans": {
        label: "简体中文",
        direction: "ltr",
      },
    },
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
        },
        blog: false,
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: "img/logo.svg",
    colorMode: {
      respectPrefersColorScheme: true,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.vsDark,
      additionalLanguages: ["toml", "powershell", "go-module"],
    },
    navbar: {
      title: "VEF Framework Go",
      logo: {
        alt: "VEF Framework Logo",
        src: "img/logo.svg",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "docsSidebar",
          position: "left",
          label: "Docs",
        },
        {
          href: "https://pkg.go.dev/github.com/coldsmirk/vef-framework-go",
          label: "API Reference",
          position: "left",
        },
        {
          type: "localeDropdown",
          position: "right",
        },
        {
          href: "https://github.com/coldsmirk/vef-framework-go",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "Introduction",
              to: "/docs/intro",
            },
            {
              label: "Getting Started",
              to: "/docs/getting-started/installation",
            },
            {
              label: "Guide",
              to: "/docs/guide/models",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "GitHub Issues",
              href: "https://github.com/coldsmirk/vef-framework-go/issues",
            },
            {
              label: "GitHub Discussions",
              href: "https://github.com/coldsmirk/vef-framework-go/discussions",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "GitHub",
              href: "https://github.com/coldsmirk/vef-framework-go",
            },
            {
              label: "Go Package",
              href: "https://pkg.go.dev/github.com/coldsmirk/vef-framework-go",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} VEF Framework Go.`,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
