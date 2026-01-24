/* TACT Web UI JavaScript */

// SSE connection status handling
document.addEventListener('htmx:sseOpen', function(evt) {
    const container = evt.target;
    container.classList.remove('sse-disconnected', 'sse-reconnecting');
});

document.addEventListener('htmx:sseError', function(evt) {
    const container = evt.target;
    container.classList.add('sse-disconnected');
    container.classList.add('sse-reconnecting');
});

document.addEventListener('htmx:sseClose', function(evt) {
    const container = evt.target;
    container.classList.add('sse-disconnected');
    container.classList.remove('sse-reconnecting');
});

// Show loading state for forms
document.addEventListener('htmx:beforeRequest', function(evt) {
    const trigger = evt.detail.elt;
    if (trigger.tagName === 'FORM') {
        const submitBtn = trigger.querySelector('button[type="submit"]');
        if (submitBtn) {
            submitBtn.setAttribute('aria-busy', 'true');
        }
    }
});

document.addEventListener('htmx:afterRequest', function(evt) {
    const trigger = evt.detail.elt;
    if (trigger.tagName === 'FORM') {
        const submitBtn = trigger.querySelector('button[type="submit"]');
        if (submitBtn) {
            submitBtn.removeAttribute('aria-busy');
        }
    }
});

// Error handling for HTMX requests
document.addEventListener('htmx:responseError', function(evt) {
    console.error('HTMX request failed:', evt.detail.xhr.status, evt.detail.xhr.statusText);
});
