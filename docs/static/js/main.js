/* ABOUTME: Client-side interactions for Gozzi docs */
/* ABOUTME: Navigation state, search, and UX enhancements */

(function() {
  'use strict';
  
  // Active link highlighting
  function highlightActiveLinks() {
    const currentPath = window.location.pathname;
    
    // Normalize paths (remove trailing slash for comparison)
    const normalizedPath = currentPath.endsWith('/') && currentPath.length > 1
      ? currentPath.slice(0, -1)
      : currentPath;
    
    // Sidebar links
    let activeLink = null;
    document.querySelectorAll('.sidebar-link').forEach(link => {
      const href = link.getAttribute('href');
      const normalizedHref = href.endsWith('/') && href.length > 1
        ? href.slice(0, -1)
        : href;
      
      if (normalizedHref === normalizedPath) {
        link.classList.add('active');
        activeLink = link;
      } else {
        link.classList.remove('active');
      }
    });
    
    // Scroll active link into view in sidebar
    if (activeLink) {
      scrollToActiveLink(activeLink);
    }
  }
  
  // Scroll active sidebar link into view
  function scrollToActiveLink(activeLink) {
    const sidebar = document.querySelector('.sidebar');
    if (!sidebar || !activeLink) return;
    
    // Use requestAnimationFrame to ensure layout is complete
    requestAnimationFrame(() => {
      const sidebarRect = sidebar.getBoundingClientRect();
      const linkRect = activeLink.getBoundingClientRect();
      
      // Check if link is not fully visible
      const isAboveView = linkRect.top < sidebarRect.top;
      const isBelowView = linkRect.bottom > sidebarRect.bottom;
      
      if (isAboveView || isBelowView) {
        // Calculate scroll position to center the link
        const sidebarScrollTop = sidebar.scrollTop;
        const linkOffsetTop = activeLink.offsetTop;
        const sidebarHeight = sidebar.clientHeight;
        const linkHeight = activeLink.clientHeight;
        
        // Center the active link in the sidebar
        const targetScroll = linkOffsetTop - (sidebarHeight / 2) + (linkHeight / 2);
        
        sidebar.scrollTo({
          top: targetScroll,
          behavior: 'smooth'
        });
      }
    });
  }
  
  // Search functionality
  function initSearch() {
    const searchBtn = document.getElementById('search-btn');
    if (searchBtn) {
      searchBtn.addEventListener('click', () => {
        alert('Search functionality coming soon! 🔍');
      });
      
      // Keyboard shortcut (Cmd+K or Ctrl+K)
      document.addEventListener('keydown', (e) => {
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
          e.preventDefault();
          searchBtn.click();
        }
      });
    }
  }
  
  // Smooth scroll for anchor links
  function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
      anchor.addEventListener('click', function(e) {
        const href = this.getAttribute('href');
        if (href === '#') return;
        
        e.preventDefault();
        const target = document.querySelector(href);
        if (target) {
          target.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
          });
          
          // Update URL without triggering scroll
          history.pushState(null, null, href);
        }
      });
    });
  }
  
  // Copy code button
  function initCodeCopy() {
    document.querySelectorAll('pre code').forEach((block) => {
      const pre = block.parentElement;
      
      // Check if button already exists
      if (pre.querySelector('.code-copy-btn')) {
        return;
      }
      
      const button = document.createElement('button');
      button.className = 'code-copy-btn';
      button.textContent = 'Copy';
      button.setAttribute('aria-label', 'Copy code to clipboard');
      
      button.addEventListener('click', async () => {
        try {
          await navigator.clipboard.writeText(block.textContent);
          button.textContent = 'Copied!';
          setTimeout(() => {
            button.textContent = 'Copy';
          }, 2000);
        } catch (err) {
          console.error('Failed to copy:', err);
          button.textContent = 'Failed';
          setTimeout(() => {
            button.textContent = 'Copy';
          }, 2000);
        }
      });
      
      pre.appendChild(button);
    });
  }
  
  // Initialize on DOM ready
  function init() {
    highlightActiveLinks();
    initSearch();
    initSmoothScroll();
    initCodeCopy();
  }
  
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
