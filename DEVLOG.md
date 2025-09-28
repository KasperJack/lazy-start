# Lazy-Start Development Log


## Development Environment

- Language: Go 1.24
- Dependencies: gorilla/mux
- Target Proxies: Nginx, Traefik, Caddy, HAProxy
- Target Services: systemctl, Docker, custom scripts





## Entry #1 - September 28, 2025

### Challenge: Request Flooding & Race Conditions
#### The Problem
When a service is starting, the loading page auto-refreshes every 2-3 seconds, causing a flood of requests to the `/start/{service}` endpoint. This creates race conditions where multiple goroutines try to start the same service simultaneously.
#### Initial Approach
```html
<meta http-equiv="refresh" content="2">
```

This caused the browser to keep hitting the same endpoint repeatedly.

#### The Realization

Actually discovered that once the service comes online, nginx naturally stops redirecting to the fallback handler. The proxy takes over and requests go directly to the real service. So the flooding stops automatically!`

#### Current Solution Direction
Moving towards WebSockets for real-time notifications:

- Client connects via WebSocket when loading page loads
- Server sends "ready" message when service is healthy
- Client redirects to original service URL
- Eliminates unnecessary polling and race conditions


### The WebSocket Reality Check
#### Challenge: WebSocket/API Endpoints Don't Work Due to Proxy Behavior

Every request to the frontend domain (eg, `jellyfin.local`) goes through this flow:

1. Request hits reverse proxy
2. Proxy tries to reach dead service → gets 502/504
3. Proxy redirects to fallback with rewrite: `/start/{service}`


No way for client to reach different endpoints than the one defined as a fallback route in the reverse proxy (path and query parameter don't work)




##### This means:

- Can't create separate `/api/status/{service}` endpoint
- Can't create WebSocket endpoints for real-time updates
- All requests get funneled through the same `/start/{service}` path 

### The Constraint-Driven Solution

The `/start/{service}` endpoint must handle ALL requests and be smart about:


- First request: Start the service + show loading page
- Subsequent requests: Just show loading page (don't restart service)
- Auto-refresh: Keep refreshing until service is up
- Natural cutoff: Once service is healthy, proxy stops hitting fallback

### Implementation Strategy

```go
func startHandler(w http.ResponseWriter, r *http.Request) {
    serviceName := mux.Vars(r)["service"]
    svc := getService(serviceName)
    
    svc.mu.Lock()
    isRunning := svc.IsRunning
    isStarting := svc.IsStarting
    svc.mu.Unlock()
    
    if isRunning {
        // Shouldn't happen, but redirect just in case
        http.Redirect(w, r, originalServiceURL, http.StatusFound)
        return
    }
    
    if !isStarting {
        // First request - start the service
        startServiceAsync(svc)
    }
    
    // Always serve loading page with auto-refresh
    serveLoadingPage(w, svc)
}

```
