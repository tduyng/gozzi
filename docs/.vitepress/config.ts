import { defineConfig } from 'vitepress'
import pkg from '../../package.json'

export default defineConfig({
    title: 'Gozzi',
    description:
        'A Go static site generator born from curiosity - Fast, flexible, and built for developers',
    lang: 'en-US',
    base: '/gozzi/', // Required for GitHub Pages on project repos (not username.github.io)

    head: [
        ['link', { rel: 'icon', type: 'image/x-icon', href: '/gozzi/favicon.ico' }],
        [
            'link',
            { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/gozzi/favicon-32x32.png' },
        ],
        [
            'link',
            { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/gozzi/favicon-16x16.png' },
        ],
        [
            'link',
            { rel: 'apple-touch-icon', sizes: '180x180', href: '/gozzi/apple-touch-icon.png' },
        ],
        ['meta', { name: 'theme-color', content: '#10b981' }],
        ['meta', { name: 'og:type', content: 'website' }],
        ['meta', { name: 'og:locale', content: 'en' }],
        ['meta', { name: 'og:site_name', content: 'Gozzi' }],
        ['meta', { name: 'og:image', content: 'https://tduyng.github.io/gozzi/logo.svg' }],
    ],

    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        logo: '/logo.svg',

        nav: [
            { text: 'Home', link: '/' },
            { text: 'Guide', link: '/guide/introduction' },
            { text: 'Reference', link: '/reference/cli/' },
            { text: 'Examples', link: '/examples/quick-start' },
            {
                text: `v${pkg.version}`,
                items: [
                    {
                        text: 'Changelog',
                        link: 'https://github.com/tduyng/gozzi/blob/main/CHANGELOG.md',
                    },
                    {
                        text: 'Contributing',
                        link: 'https://github.com/tduyng/gozzi/blob/main/CONTRIBUTING.md',
                    },
                ],
            },
        ],

        sidebar: {
            '/guide/': [
                {
                    text: 'Getting Started',
                    items: [
                        { text: 'What is Gozzi?', link: '/guide/introduction' },
                        { text: 'Installation', link: '/guide/installation' },
                        { text: 'Your First Site', link: '/guide/getting-started' },
                    ],
                },
                {
                    text: 'Core Concepts',
                    items: [
                        { text: 'Content Structure', link: '/guide/concepts/content-structure' },
                    ],
                },
                {
                    text: 'Configuration',
                    collapsed: false,
                    items: [
                        { text: 'Overview', link: '/guide/configuration/' },
                        { text: 'Site Settings', link: '/guide/configuration/site' },
                        { text: 'Environment Configs', link: '/guide/configuration/environment' },
                        { text: 'Extended Config', link: '/guide/configuration/extended' },
                        { text: 'Front Matter', link: '/guide/configuration/frontmatter' },
                        { text: 'Sections', link: '/guide/configuration/sections' },
                        { text: 'Pages', link: '/guide/configuration/pages' },
                        { text: 'Inheritance', link: '/guide/configuration/inheritance' },
                        { text: 'Troubleshooting', link: '/guide/configuration/troubleshooting' },
                    ],
                },
                {
                    text: 'Templates',
                    collapsed: false,
                    items: [
                        { text: 'Overview', link: '/guide/templates/' },
                        { text: 'Structure', link: '/guide/templates/structure' },
                        { text: 'Template Mapping', link: '/guide/templates/mapping' },
                        { text: 'Development', link: '/guide/templates/development' },
                        { text: 'Variables', link: '/guide/templates/variables' },
                        { text: 'Functions', link: '/guide/templates/functions' },
                        { text: 'Inheritance', link: '/guide/templates/inheritance' },
                        { text: 'Partials', link: '/guide/templates/partials' },
                        { text: 'Macros', link: '/guide/templates/macros' },
                        { text: 'Advanced Patterns', link: '/guide/templates/advanced' },
                        { text: 'Examples', link: '/guide/templates/examples' },
                    ],
                },
                {
                    text: 'Built-in Features',
                    collapsed: false,
                    items: [
                        { text: 'Overview', link: '/guide/features/' },
                        { text: 'Mathematical Expressions', link: '/guide/features/math' },
                        { text: 'Mermaid Diagrams', link: '/guide/features/diagrams' },
                        {
                            text: 'Syntax Highlighting',
                            link: '/guide/features/syntax-highlighting',
                        },
                        { text: 'Content Features', link: '/guide/features/content-features' },
                        { text: 'SEO Automation', link: '/guide/features/seo' },
                    ],
                },
            ],
            '/reference/': [
                {
                    text: 'CLI Reference',
                    collapsed: false,
                    items: [
                        { text: 'Overview', link: '/reference/cli/' },
                        { text: 'Build Command', link: '/reference/cli/build' },
                        { text: 'Serve Command', link: '/reference/cli/serve' },
                        { text: 'Usage Patterns', link: '/reference/cli/usage-patterns' },
                        { text: 'Troubleshooting', link: '/reference/cli/troubleshooting' },
                    ],
                },
                {
                    text: 'Reference',
                    items: [
                        { text: 'Template Functions', link: '/reference/template-functions' },
                        { text: 'Architecture', link: '/reference/architecture' },
                    ],
                },
            ],
            '/examples/': [
                {
                    text: 'Examples',
                    items: [
                        { text: 'Complete Blog Tutorial', link: '/examples/quick-start' },
                        { text: 'Real-World: tduyng.com', link: '/examples/real-world' },
                    ],
                },
            ],
        },

        socialLinks: [{ icon: 'github', link: 'https://github.com/tduyng/gozzi' }],

        search: {
            provider: 'local',
        },

        editLink: {
            pattern: 'https://github.com/tduyng/gozzi/edit/main/docs/:path',
            text: 'Edit this page on GitHub',
        },

        footer: {
            message: 'Released under the MIT License.',
            copyright: 'Copyright © 2024-present Duy NG',
        },

        outline: {
            level: [2, 3],
            label: 'On this page',
        },
    },

    markdown: {
        theme: {
            light: 'github-light',
            dark: 'github-dark',
        },
        lineNumbers: true,
    },

    ignoreDeadLinks: [/^http:\/\/localhost/, './README'],
})
