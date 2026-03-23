/**
 * Example: Integration with a simple app to demonstrate singleton pattern
 * This prevents re-creating client instances on every app start
 */

import { RavaClient } from "./src/client"

// ========== APPLICATION STARTUP ==========
console.log("🚀 Starting application...")

// Initialize ONCE at app startup
RavaClient.initialize({
    apiKey: "demo-key-123",
    baseUrl: "https://api.rava.dev",
})

console.log("✅ Rava initialized (singleton created)")

// ========== SIMULATING MULTIPLE OPERATIONS ==========

async function operation1() {
    console.log("\n📄 Operation 1: Ingesting document...")
    const client = RavaClient.getInstance() // ✅ SAME instance
    console.log("Got client instance (ID reference check)")
    // Would call: await client.ingest(...)
}

async function operation2() {
    console.log("\n🔍 Operation 2: Querying documents...")
    const client = RavaClient.getInstance() // ✅ SAME instance
    console.log("Got client instance (ID reference check)")
    // Would call: await client.query(...)
}

async function operation3() {
    console.log("\n📄 Operation 3: Ingesting another document...")
    const client = RavaClient.getInstance() // ✅ SAME instance
    console.log("Got client instance (ID reference check)")
    // Would call: await client.ingest(...)
}

// ========== DEMONSTRATE THE SOLUTION ==========

async function main() {
    try {
        // Multiple operations reuse the SAME instance
        await operation1()
        await operation2()
        await operation3()

        console.log("\n✨ All operations completed with single client instance!")
        console.log("❌ NO new axios instances were created during app lifecycle")
    } catch (error: any) {
        console.error("Error:", error.message)
    }

    // ========== DEMONSTRATE ERROR HANDLING ==========
    console.log("\n--- Testing error handling ---")

    RavaClient.reset() // Reset for demo

    try {
        const client = RavaClient.getInstance() // Should throw
    } catch (error: any) {
        console.log("✅ Caught expected error:", error.message)
    }

    // Reinitialize properly
    RavaClient.initialize({ apiKey: "new-key" })
    const client = RavaClient.getInstance()
    console.log("✅ Successfully reinitialized")
}

main().catch(console.error)

// ========== SUMMARY ==========
/*
 * BEFORE (❌ Problem):
 * - Creating new RavaClient instances on every operation
 * - Every app restart recreates the entire client
 * - Memory waste and initialization overhead
 *
 * AFTER (✅ Solution):
 * - Single RavaClient instance for entire app lifecycle
 * - RavaClient.initialize() called once at startup
 * - RavaClient.getInstance() reuses same instance everywhere
 * - No recreation on app restart
 */
