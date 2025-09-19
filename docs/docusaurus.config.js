const { themes } = require('prism-react-renderer');
const lightCodeTheme = themes.github;
const darkCodeTheme = themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'FoundryVTT REST API Relay',
  tagline: 'Message relay server for the foundryvtt-rest-api Foundry Module',
  favicon: 'img/favicon.svg',

  // Set the production url of your site here
  url: 'https://foundryvtt-rest-api-relay.fly.dev',
  // Set the /<baseUrl>/ pathname under which your site is served
  baseUrl: '/docs/',

  // Add trailing slash configuration
  trailingSlash: false,

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          path: './md',
          routeBasePath: '/',
          sidebarPath: './sidebars.js',
        },
        blog: false, // Disable blog
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  plugins: [
    // This plugin will generate the API documentation from the TypeScript files
    // We're using our custom generateApiDocs.js script instead of TypeDoc
    // because we want to include inline comments and createApiRoute configuration
    async function customDocPlugin(context, options) {
      return {
        name: 'custom-api-docs-plugin',
        async loadContent() {
          // Execute our custom generator script during build
          try {
            const generateDocs = require('../scripts/generateApiDocs.js');
            // Script execution is handled in the file itself
          } catch (error) {
            console.error('Error generating API documentation:', error);
          }
          return null;
        },
      };
    },
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.jpg',
      navbar: {
        title: 'FoundryVTT REST API Relay',
        logo: {
          alt: 'FoundryVTT Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'tutorialSidebar',
            position: 'left',
            label: 'Documentation',
          },
          {
            to: '/api/',
            label: 'API Reference',
            position: 'left',
          },
          {
            href: 'https://github.com/ThreeHats/foundryvtt-rest-api-relay',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Getting Started',
                to: '/intro',
              },
              {
                label: 'API Reference',
                to: '/api/',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Discord',
                href: 'https://discord.gg/U634xNGRAC',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/ThreeHats/foundryvtt-rest-api-relay',
              },
            ],
          },
        ],
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      colorMode: {
        defaultMode: 'dark',
        disableSwitch: false,
        respectPrefersColorScheme: false,
      },
      algolia: {
        // The application ID provided by Algolia
        appId: 'YOUR_APP_ID',
        
        // Public API key: it is safe to commit it
        apiKey: 'YOUR_SEARCH_API_KEY',
        
        indexName: 'foundryvtt-rest-api-relay',
        
        // Optional: see doc section below
        contextualSearch: true,
        
        // Optional: Specify domains where the navigation should occur through window.location instead on history.push. Useful when our Algolia config crawls multiple documentation sites and we want to navigate with window.location.href to them.
        externalUrlRegex: 'external\\.com|domain\\.com',
        
        // Optional: Replace parts of the item URLs from Algolia. Useful when using the same search index for multiple deployments using a different baseUrl. You can use regexp or string in the `from` param. For example: localhost:3000 vs myCompany.github.io/myDocusaurusWebsite/
        replaceSearchResultPathname: {
          from: '/docs/', // or as RegExp: /\/docs\//
          to: '/',
        },
        
        // Optional: Algolia search parameters
        searchParameters: {},
        
        // Optional: path for search page that enabled by default (`false` to disable it)
        searchPagePath: 'search',
        
        // ... other Algolia params
      },
    }),
};

export default config;
