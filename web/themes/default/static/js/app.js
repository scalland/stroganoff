// GOCR Web Interface - Application JavaScript

document.addEventListener('DOMContentLoaded', function() {
    console.log('GOCR application initialized');

    // Add event listeners for navigation
    setupNavigation();

    // Initialize API calls
    initializeAPI();
});

function setupNavigation() {
    const links = document.querySelectorAll('.navbar-menu a');

    links.forEach(link => {
        link.addEventListener('click', function(e) {
            const href = this.getAttribute('href');

            if (href.startsWith('#')) {
                e.preventDefault();
                const target = document.querySelector(href);

                if (target) {
                    target.scrollIntoView({ behavior: 'smooth' });
                }
            }
        });
    });
}

function initializeAPI() {
    // Check health status
    checkHealth();

    // Check heartbeat every 30 seconds
    setInterval(checkHeartbeat, 30000);
}

function checkHealth() {
    fetch('/health')
        .then(response => response.json())
        .then(data => {
            console.log('Health check passed:', data);
        })
        .catch(error => {
            console.error('Health check failed:', error);
        });
}

function checkHeartbeat() {
    fetch('/api/heartbeat')
        .then(response => response.json())
        .then(data => {
            console.log('Heartbeat received:', data);
        })
        .catch(error => {
            console.error('Heartbeat failed:', error);
        });
}

// API helper functions
async function getMetrics(token) {
    try {
        const response = await fetch('/api/metrics', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        return await response.json();
    } catch (error) {
        console.error('Failed to fetch metrics:', error);
    }
}

async function createToken(scopes = [], duration = 86400) {
    try {
        const response = await fetch('/api/auth/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                scopes: scopes,
                duration: duration
            })
        });
        return await response.json();
    } catch (error) {
        console.error('Failed to create token:', error);
    }
}
