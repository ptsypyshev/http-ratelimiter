package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// client struct should be used to decide to limit the connection or not
type client struct {
	limiter      *rate.Limiter // limiter instance that controls rate limit
	blocked      bool          // blocked status of client
	blockedSince time.Time     // the time when client was blocked
	lastSeen     time.Time     // the time when client was seen last time
}

// rateLimiter is a middleware module that can limit the connection from clients
func (a *App) rateLimiter(c *gin.Context) {
	clientIP := getIP(c)                                  // get client IP
	clientNetwork, err := getNetwork(clientIP, a.netmask) // get client network
	if err != nil {
		c.String(http.StatusBadRequest, "%s\n", err) // Send 400 Status code if get an invalid network
		c.Abort()                                    // Stop next processing steps
	}

	a.mu.Lock() // Lock access to clientsRate map

	// Create new key-value pair in clientsRate map if no client network key is found
	if _, found := a.clientsRate[clientNetwork]; !found {
		a.clientsRate[clientNetwork] = &client{
			limiter:  rate.NewLimiter(rate.Every(time.Minute/time.Duration(a.limitPerMinute)), a.limitPerMinute),
			lastSeen: time.Now(),
		}
	}

	// Check blocked status and cooldown period
	if a.clientsRate[clientNetwork].blocked && time.Since(a.clientsRate[clientNetwork].blockedSince) < a.cooldownPeriod {
		a.mu.Unlock()                                               // Unlock access to clientsRate map
		c.String(http.StatusTooManyRequests, "Too many requests\n") // Send 429 Status code and simple string as a page content
		c.Abort()                                                   // Stop next processing steps
		return
	}

	// Check the limiter status
	if !a.clientsRate[clientNetwork].limiter.Allow() {
		a.clientsRate[clientNetwork].blocked = true            // Set blocked status to client network
		a.clientsRate[clientNetwork].blockedSince = time.Now() // update blocked since time
		a.mu.Unlock()                                          // Unlock access to clientsRate map

		c.String(http.StatusTooManyRequests, "Too many requests\n") // Send 429 Status code and simple string as a page content
		c.Abort()                                                   // Stop next processing steps
		return
	} else {
		a.clientsRate[clientNetwork].blocked = false // Set unblocked status to client network
	}

	a.mu.Unlock() // Unlock access to clientsRate map

	c.Next() // Execute next step handler
}

// getIP returns string representation of client IP address
// from "X-Forwarded-For" header field or from Request.RemoteAddr (if "X-Forwarded-For" is not filled).
func getIP(c *gin.Context) string {
	clientIP := c.Request.Header.Get("X-Forwarded-For") // get a value from "X-Forwarded-For" header field
	if clientIP == "" {
		clientIP = c.RemoteIP() // get a value from Request.RemoteAddr
	}
	return clientIP
}

// getNetwork returns a string representation of client network address
// or error if the CIDR cannot be parsed.
func getNetwork(ip string, netmask uint8) (string, error) {
	CIDR := fmt.Sprintf("%s/%d", ip, netmask) // create CIDR string from client IP address and netmask value
	_, ipnet, err := net.ParseCIDR(CIDR)      // Try to parse CIDR
	if err != nil {
		return "", errors.New("bad client IP address received")
	}

	clientNetwork := ipnet.String()[:len(ipnet.String())-3] // get network string representation from ipnet instance
	return clientNetwork, nil
}
