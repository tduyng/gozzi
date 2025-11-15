import { defineConfig } from 'vitepress'
import { readFileSync } from 'fs'
import { resolve } from 'path'

// Read version from VERSION file
const version = readFileSync(resolve(__dirname, '../../VERSION'), 'utf-8').trim()

// https://vitepress.dev/reference/site-config
export default defineConfig({
    title: 'Gozzi',
    description:
        'A Go static site generator born from curiosity - Fast, flexible, and built for developers',
    lang: 'en-US',

    head: [
        ['link', { rel: 'icon', href: '/favicon.ico' }],
        ['meta', { name: 'theme-color', content: '#10b981' }],
        ['meta', { name: 'og:type', content: 'website' }],
        ['meta', { name: 'og:locale', content: 'en' }],
        ['meta', { name: 'og:site_name', content: 'Gozzi' }],
        ['meta', { name: 'og:image', content: 'https://gozzi.dev/og-image.png' }],
    ],

    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        logo: '/logo.svg',

        nav: [
            { text: 'Home', link: '/' },
            { text: 'Guide', link: '/guide/introduction' },
            { text: 'Reference', link: '/reference/cli' },
            { text: 'Examples', link: '/examples/quick-start' },
            {
                text: `v${version}`,
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
                    text: 'Guides',
                    items: [
                        { text: 'Content Structure', link: '/guide/content-structure' },
                        { text: 'Configuration', link: '/guide/configuration' },
                        { text: 'Templates', link: '/guide/templates' },
                        { text: 'Built-in Features', link: '/guide/features' },
                    ],
                },
            ],
            '/reference/': [
                {
                    text: 'Reference',
                    items: [
                        { text: 'CLI Commands', link: '/reference/cli' },
                        { text: 'Template Functions', link: '/reference/template-functions' },
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

    ignoreDeadLinks: [
        /^http:\/\/localhost/,
        './README'
    ]
})
