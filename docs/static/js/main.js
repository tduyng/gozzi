;(function () {
    'use strict'

    // Platform detection
    const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0

    // Mobile drawer
    function initDrawer() {
        const menuToggle = document.getElementById('menu-toggle')
        const closeDrawer = document.getElementById('close-drawer')
        const drawer = document.getElementById('sidebar-drawer')
        const overlay = document.getElementById('sidebar-overlay')

        if (!menuToggle || !drawer) return

        function openDrawer() {
            drawer.classList.add('open')
            overlay?.classList.add('open')
            document.body.classList.add('drawer-open')
            menuToggle.setAttribute('aria-expanded', 'true')
            drawer.setAttribute('aria-hidden', 'false')
            closeDrawer?.focus()
        }

        function closeDrawerFn() {
            drawer.classList.remove('open')
            overlay?.classList.remove('open')
            document.body.classList.remove('drawer-open')
            menuToggle.setAttribute('aria-expanded', 'false')
            drawer.setAttribute('aria-hidden', 'true')
            menuToggle.focus()
        }

        menuToggle.addEventListener('click', openDrawer)
        closeDrawer?.addEventListener('click', closeDrawerFn)
        overlay?.addEventListener('click', closeDrawerFn)

        // Close on escape
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && drawer.classList.contains('open')) {
                closeDrawerFn()
            }
        })

        // Sync sidebar links with drawer
        const sidebarLinks = document.querySelectorAll('.sidebar-link')
        const drawerLinks = document.querySelectorAll('.sidebar-drawer .sidebar-link')
        
        function syncActiveLink() {
            const currentPath = window.location.pathname
            const normalizedPath = currentPath.endsWith('/') && currentPath.length > 1
                ? currentPath.slice(0, -1)
                : currentPath

            sidebarLinks.forEach(link => {
                const href = link.getAttribute('href')
                const normalizedHref = href?.endsWith('/') && href.length > 1 ? href.slice(0, -1) : href
                if (normalizedHref === normalizedPath) {
                    link.classList.add('active')
                } else {
                    link.classList.remove('active')
                }
            })

            drawerLinks.forEach(link => {
                const href = link.getAttribute('href')
                const normalizedHref = href?.endsWith('/') && href.length > 1 ? href.slice(0, -1) : href
                if (normalizedHref === normalizedPath) {
                    link.classList.add('active')
                } else {
                    link.classList.remove('active')
                }
            })
        }

        syncActiveLink()
    }

    // Active link highlighting
    function highlightActiveLinks() {
        const currentPath = window.location.pathname
        const normalizedPath = currentPath.endsWith('/') && currentPath.length > 1
            ? currentPath.slice(0, -1)
            : currentPath

        document.querySelectorAll('.sidebar-link').forEach((link) => {
            const href = link.getAttribute('href')
            const normalizedHref = href?.endsWith('/') && href.length > 1 ? href.slice(0, -1) : href

            if (normalizedHref === normalizedPath) {
                link.classList.add('active')
            } else {
                link.classList.remove('active')
            }
        })
    }

    // Scroll active sidebar link into view
    function scrollToActiveLink() {
        const activeLink = document.querySelector('.sidebar-link.active')
        const sidebar = document.querySelector('.sidebar')
        if (!sidebar || !activeLink) return

        requestAnimationFrame(() => {
            const sidebarRect = sidebar.getBoundingClientRect()
            const linkRect = activeLink.getBoundingClientRect()

            const isAboveView = linkRect.top < sidebarRect.top
            const isBelowView = linkRect.bottom > sidebarRect.bottom

            if (isAboveView || isBelowView) {
                const sidebarScrollTop = sidebar.scrollTop
                const linkOffsetTop = activeLink.offsetTop
                const sidebarHeight = sidebar.clientHeight
                const linkHeight = activeLink.clientHeight

                const targetScroll = linkOffsetTop - sidebarHeight / 2 + linkHeight / 2

                sidebar.scrollTo({
                    top: targetScroll,
                    behavior: 'smooth',
                })
            }
        })
    }

    // Search functionality
    function initSearch() {
        const searchBtn = document.getElementById('search-btn')
        const searchKbd = searchBtn?.querySelector('.search-kbd')
        
        // Update keyboard shortcut based on platform
        if (searchKbd) {
            searchKbd.textContent = isMac ? '⌘K' : 'Ctrl+K'
        }

        if (searchBtn) {
            searchBtn.addEventListener('click', () => {
                // Placeholder for search - could integrate Algolia, Fuse.js, etc.
                console.log('Search triggered')
            })

            // Keyboard shortcut
            document.addEventListener('keydown', (e) => {
                const key = isMac ? e.metaKey : e.ctrlKey
                if (key && e.key === 'k') {
                    e.preventDefault()
                    searchBtn.click()
                }
            })
        }
    }

    // Smooth scroll for anchor links
    function initSmoothScroll() {
        document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
            anchor.addEventListener('click', function (e) {
                const href = this.getAttribute('href')
                if (href === '#') return

                e.preventDefault()
                const target = document.querySelector(href)
                if (target) {
                    target.scrollIntoView({
                        behavior: 'smooth',
                        block: 'start',
                    })
                    history.pushState(null, null, href)
                }
            })
        })
    }

    // Copy code button
    function initCodeCopy() {
        document.querySelectorAll('pre code').forEach((block) => {
            const pre = block.parentElement

            if (pre.querySelector('.code-copy-btn')) {
                return
            }

            const button = document.createElement('button')
            button.className = 'code-copy-btn'
            button.textContent = 'Copy'
            button.setAttribute('aria-label', 'Copy code to clipboard')

            button.addEventListener('click', async () => {
                try {
                    await navigator.clipboard.writeText(block.textContent)
                    button.textContent = 'Copied!'
                    setTimeout(() => {
                        button.textContent = 'Copy'
                    }, 2000)
                } catch (err) {
                    console.error('Failed to copy:', err)
                    button.textContent = 'Failed'
                    setTimeout(() => {
                        button.textContent = 'Copy'
                    }, 2000)
                }
            })

            pre.appendChild(button)
        })
    }

    // Initialize
    function init() {
        initDrawer()
        highlightActiveLinks()
        scrollToActiveLink()
        initSearch()
        initSmoothScroll()
        initCodeCopy()
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init)
    } else {
        init()
    }
})()
