;(function () {
    'use strict'

    const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0

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

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && drawer.classList.contains('open')) {
                closeDrawerFn()
            }
        })

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

    function initSearch() {
        const searchBtn = document.getElementById('search-btn')
        const searchKbd = searchBtn?.querySelector('.search-kbd')
        
        if (searchKbd) {
            searchKbd.textContent = isMac ? '⌘K' : 'Ctrl+K'
        }

        if (searchBtn) {
            searchBtn.addEventListener('click', () => {
                console.log('Search triggered')
            })

            document.addEventListener('keydown', (e) => {
                const key = isMac ? e.metaKey : e.ctrlKey
                if (key && e.key === 'k') {
                    e.preventDefault()
                    searchBtn.click()
                }
            })
        }
    }

    function initSmoothScroll() {
        document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
            anchor.addEventListener('click', function (e) {
                const href = this.getAttribute('href')
                if (!href || href === '#' || !href.startsWith('#')) return

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

    function initCodeCopy() {
        document.querySelectorAll('pre code').forEach((block) => {
            const pre = block.parentElement

            if (pre.querySelector('.code-copy-btn')) {
                return
            }

            const button = document.createElement('button')
            button.className = 'code-copy-btn'
            button.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg> Copy'
            button.setAttribute('aria-label', 'Copy code to clipboard')

            button.addEventListener('click', async () => {
                try {
                    await navigator.clipboard.writeText(block.textContent)
                    button.classList.add('copied')
                    button.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg> Copied!'
                    setTimeout(() => {
                        button.classList.remove('copied')
                        button.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg> Copy'
                    }, 2000)
                } catch (err) {
                    console.error('Failed to copy:', err)
                    button.innerHTML = 'Failed'
                    setTimeout(() => {
                        button.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg> Copy'
                    }, 2000)
                }
            })

            pre.appendChild(button)
        })
    }

    function initHeaderScroll() {
        const header = document.querySelector('.header')
        if (!header) return

        let lastScroll = 0
        const scrollThreshold = 10

        function handleScroll() {
            const currentScroll = window.pageYOffset || document.documentElement.scrollTop

            if (Math.abs(currentScroll - lastScroll) < scrollThreshold) {
                return
            }

            if (currentScroll > 50) {
                header.classList.add('scrolled')
            } else {
                header.classList.remove('scrolled')
            }

            lastScroll = currentScroll <= 0 ? 0 : currentScroll
        }

        window.addEventListener('scroll', handleScroll, { passive: true })
        handleScroll()
    }

    function initCollapsibleSidebar() {
        const sidebarSections = document.querySelectorAll('.sidebar-section.collapsible')

        sidebarSections.forEach((section) => {
            const title = section.querySelector('.sidebar-title')
            if (!title) return

            const storedState = localStorage.getItem(`sidebar-${section.dataset.section}`)
            if (storedState === 'collapsed') {
                section.classList.add('collapsed')
            }

            title.addEventListener('click', () => {
                section.classList.toggle('collapsed')
                const isCollapsed = section.classList.contains('collapsed')
                localStorage.setItem(`sidebar-${section.dataset.section}`, isCollapsed ? 'collapsed' : 'expanded')
            })
        })
    }

    function initAnchorLinks() {
        const prose = document.querySelector('.prose')
        if (!prose) return

        const headings = prose.querySelectorAll('h2, h3, h4')

        headings.forEach((heading) => {
            if (!heading.id) {
                const text = heading.textContent.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
                heading.id = text
            }

            const anchorLink = document.createElement('a')
            anchorLink.className = 'anchor-link'
            anchorLink.href = `#${heading.id}`
            anchorLink.setAttribute('aria-hidden', 'true')
            anchorLink.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path></svg>'

            anchorLink.addEventListener('click', (e) => {
                e.preventDefault()
                const target = document.getElementById(heading.id)
                if (target) {
                    target.scrollIntoView({ behavior: 'smooth', block: 'start' })
                    history.pushState(null, null, `#${heading.id}`)
                }
            })

            heading.style.position = 'relative'
            heading.insertBefore(anchorLink, heading.firstChild)
        })
    }

    function initTableOfContents() {
        const toc = document.querySelector('.toc')
        const prose = document.querySelector('.prose')
        
        if (!toc || !prose) return

        const headings = prose.querySelectorAll('h2, h3')
        const tocList = toc.querySelector('.toc-list')
        
        if (!tocList || headings.length === 0) {
            toc.style.display = 'none'
            return
        }

        const fragment = document.createDocumentFragment()

        headings.forEach((heading) => {
            const text = heading.textContent
            let id = heading.id
            
            if (!id) {
                id = text.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
                heading.id = id
            }

            const listItem = document.createElement('li')
            const link = document.createElement('a')
            link.href = `#${id}`
            link.textContent = text
            
            if (heading.tagName === 'H3') {
                link.classList.add('toc-h3')
            }

            link.addEventListener('click', (e) => {
                e.preventDefault()
                heading.scrollIntoView({ behavior: 'smooth', block: 'start' })
                history.pushState(null, null, `#${id}`)
            })

            listItem.appendChild(link)
            fragment.appendChild(listItem)
        })

        tocList.innerHTML = ''
        tocList.appendChild(fragment)

        const observerOptions = {
            rootMargin: '-80px 0px -80% 0px',
            threshold: 0
        }

        const observer = new IntersectionObserver((entries) => {
            entries.forEach((entry) => {
                if (entry.isIntersecting) {
                    const id = entry.target.id
                    tocList.querySelectorAll('a').forEach((link) => {
                        link.classList.remove('active')
                        if (link.getAttribute('href') === `#${id}`) {
                            link.classList.add('active')
                        }
                    })
                }
            })
        }, observerOptions)

        headings.forEach((heading) => observer.observe(heading))
    }

    function initBreadcrumbs() {
        const breadcrumbsContainer = document.querySelector('.breadcrumbs')
        if (!breadcrumbsContainer) return

        const path = window.location.pathname
        const basePath = window.location.origin
        const pathParts = path.split('/').filter(Boolean)

        if (pathParts.length === 0 || (pathParts.length === 1 && pathParts[0] === '')) {
            breadcrumbsContainer.style.display = 'none'
            return
        }

        const breadcrumbs = []
        
        const homeItem = document.createElement('a')
        homeItem.href = '/'
        homeItem.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>'
        homeItem.setAttribute('aria-label', 'Home')
        breadcrumbs.push(homeItem)

        let currentPath = ''
        pathParts.forEach((part, index) => {
            currentPath += `/${part}`

            const isLast = index === pathParts.length - 1
            
            const separator = document.createElement('span')
            separator.className = 'breadcrumbs-separator'
            separator.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>'
            breadcrumbs.push(separator)

            if (isLast) {
                const currentItem = document.createElement('span')
                currentItem.className = 'breadcrumbs-current'
                currentItem.textContent = part.replace(/-/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())
                breadcrumbs.push(currentItem)
            } else {
                const linkItem = document.createElement('a')
                linkItem.href = currentPath
                linkItem.textContent = part.replace(/-/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())
                breadcrumbs.push(linkItem)
            }
        })

        breadcrumbsContainer.innerHTML = ''
        breadcrumbs.forEach((item) => breadcrumbsContainer.appendChild(item))
    }

    function initPrevNext() {
        const prevNextContainer = document.querySelector('.prev-next')
        const sidebarLinks = Array.from(document.querySelectorAll('.sidebar-link'))
        
        if (!prevNextContainer || sidebarLinks.length < 2) return

        const currentPath = window.location.pathname
        const normalizedPath = currentPath.endsWith('/') && currentPath.length > 1
            ? currentPath.slice(0, -1)
            : currentPath

        let currentIndex = -1

        sidebarLinks.forEach((link, index) => {
            const href = link.getAttribute('href')
            const normalizedHref = href?.endsWith('/') && href.length > 1 ? href.slice(0, -1) : href
            if (normalizedHref === normalizedPath) {
                currentIndex = index
            }
        })

        if (currentIndex === -1) return

        const prevLink = currentIndex > 0 ? sidebarLinks[currentIndex - 1] : null
        const nextLink = currentIndex < sidebarLinks.length - 1 ? sidebarLinks[currentIndex + 1] : null

        const prevElement = prevNextContainer.querySelector('.prev-next-link.prev')
        const nextElement = prevNextContainer.querySelector('.prev-next-link.next')

        if (prevElement && prevLink) {
            const prevTitle = prevElement.querySelector('.prev-next-title')
            if (prevTitle) {
                prevTitle.textContent = prevLink.textContent
            }
            prevElement.href = prevLink.getAttribute('href')
            prevElement.style.display = 'flex'
        } else if (prevElement) {
            prevElement.style.display = 'none'
        }

        if (nextElement && nextLink) {
            const nextTitle = nextElement.querySelector('.prev-next-title')
            if (nextTitle) {
                nextTitle.textContent = nextLink.textContent
            }
            nextElement.href = nextLink.getAttribute('href')
            nextElement.style.display = 'flex'
        } else if (nextElement) {
            nextElement.style.display = 'none'
        }
    }

    function initTabs() {
        const tabContainers = document.querySelectorAll('.tabs')

        tabContainers.forEach((container) => {
            const tabs = container.querySelectorAll('.tabs-tab')
            const panels = container.querySelectorAll('.tabs-panel')

            tabs.forEach((tab) => {
                tab.addEventListener('click', () => {
                    const target = tab.dataset.tab

                    tabs.forEach((t) => t.classList.remove('active'))
                    panels.forEach((p) => p.classList.remove('active'))

                    tab.classList.add('active')
                    const panel = container.querySelector(`#${target}`)
                    if (panel) {
                        panel.classList.add('active')
                    }
                })
            })
        })
    }

    function initAccordion() {
        const accordionItems = document.querySelectorAll('.accordion-item')

        accordionItems.forEach((item) => {
            const header = item.querySelector('.accordion-header')
            if (!header) return

            header.addEventListener('click', () => {
                const isOpen = item.classList.contains('open')
                
                accordionItems.forEach((i) => i.classList.remove('open'))

                if (!isOpen) {
                    item.classList.add('open')
                }
            })
        })
    }

    function initExternalLinks() {
        const prose = document.querySelector('.prose')
        if (!prose) return

        const links = prose.querySelectorAll('a')
        links.forEach((link) => {
            const href = link.getAttribute('href')
            if (href && href.startsWith('http') && !href.includes(window.location.hostname)) {
                link.classList.add('external-link')
                link.setAttribute('target', '_blank')
                link.setAttribute('rel', 'noopener noreferrer')
            }
        })
    }

    function init() {
        initDrawer()
        highlightActiveLinks()
        scrollToActiveLink()
        initSearch()
        initSmoothScroll()
        initCodeCopy()
        initHeaderScroll()
        initCollapsibleSidebar()
        initAnchorLinks()
        initTableOfContents()
        initBreadcrumbs()
        initPrevNext()
        initTabs()
        initAccordion()
        initExternalLinks()
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init)
    } else {
        init()
    }
})()
