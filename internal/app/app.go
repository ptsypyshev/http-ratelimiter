// Package app provides web application that can handle HTTP requests
package app

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimitPerMinute          = 100 // Default limit of requests per minute
	defaultCooldownPeriodInMinutes = 1   // Default cooldown period (application will response with 429 status code)
	defaultCleanPeriodInMinutes    = 1   // Defaulc clean period for App.clientsRate map
	defaultBanNetmask              = 24  // Default network mask for BAN
)

// App presents a http rate limiter functionality for web application
type App struct {
	router         *gin.Engine        // http router (used Gin framework)
	mu             *sync.Mutex        // Mutex that used to prevent race during async operations
	clientsRate    map[string]*client // Map which uses ip-network as a key and *client struct as a value
	limitPerMinute int                // limit of requests per minute
	cooldownPeriod time.Duration      // cooldown period (application will response with 429 status code)
	cleanPeriod    time.Duration      // clean period (separate goroutine cleans the expired key-value pairs from clientsRate map)
	netmask        uint8              // network mask used to get network address from client IP address
}

// Constructor for App instance
func NewApp() *App {
	limitPerMinute, netmask, cooldownTime, cleanPeriod := getVars() // Gets variables from Environment or default values
	return &App{
		router:         gin.New(),
		mu:             &sync.Mutex{},
		clientsRate:    make(map[string]*client),
		limitPerMinute: limitPerMinute,
		cooldownPeriod: cooldownTime,
		cleanPeriod:    cleanPeriod,
		netmask:        netmask,
	}
}

// Run method used to
func (a *App) Run() {
	a.router.Use(a.rateLimiter)  // Plug in the rate limiter middleware module
	a.router.Use(gin.Recovery()) // Plug in default recovery mechanism for Gin framework

	// Set routes
	a.router.GET("/", a.indexHandler)
	a.router.GET("/clear/", a.clearLimitsHandler)

	// Run in background function that cleans clientsRate map
	go a.cleanClientsRateMap()

	// Run web-server
	if err := a.router.Run(":8080"); err != nil {
		log.Fatalf("%s\n", err)
	}
}

// cleanClientsRateMap deletes the expired key-value pairs from clientsRate map
func (a *App) cleanClientsRateMap() {
	// Run infinity loop
	for {
		a.mu.Lock() // Lock access to clientsRate map

		// Loop the range of clientsRate map to delete the expired key-value pairs
		for key, client := range a.clientsRate {
			if time.Since(client.lastSeen) > a.cooldownPeriod {
				delete(a.clientsRate, key)
			}
		}

		a.mu.Unlock()             // Unlock access to clientsRate map
		time.Sleep(a.cleanPeriod) // Sleep for clean period duration
	}
}

// getVars returns values of limitPerMinute, cooldownPeriod, cleanPeriod and netmask from OS Environment or from default values
func getVars() (limitPerMinute int, netmask uint8, cooldownPeriod, cleanPeriod time.Duration) {
	limitPerMinute = getIntEnv("LIMIT_PER_MINUTE", defaultLimitPerMinute)
	cooldownPeriodVal := getIntEnv("COOLDOWN_PERIOD_IN_MINUTES", defaultCooldownPeriodInMinutes)
	cooldownPeriod = time.Duration(cooldownPeriodVal) * time.Minute
	cleanPeriodVal := getIntEnv("CLEAN_PERIOD_IN_MINUTES", defaultCleanPeriodInMinutes)
	cleanPeriod = time.Duration(cleanPeriodVal) * time.Minute
	netmask = uint8(getIntEnv("NETMASK", defaultBanNetmask))
	return
}

// getIntEnv returns integer value of OS Environment variable or it's default value if something went wrong
func getIntEnv(envName string, defaultValue int) int {
	env, ok := os.LookupEnv(envName) // Looks up variable at OS Environment
	if !ok {
		return defaultValue // OS Environment variable was not set
	}
	result, err := strconv.Atoi(env) // Try to convert variable from string to int
	if err != nil {
		return defaultValue // Cannot convert from string to int and returns default value
	}
	return result // Returns converted value
}
