#!/bin/bash

echo "Generating status codes..."

# Define status codes
declare -A status_codes
status_codes=(
    ["GameFound"] = 10
	["GameStarted"] = 11
	["GameCountdown"] = 12
	["PlayerJoined"] = 20
	["PlayerLeft"] = 30
    )

# Generate TypeScript code
echo "export enum StatusCode {" > statusCodes.ts
for name in "${!status_codes[@]}"; do
    echo "    $name = ${status_codes[$name]}," >> statusCodes.ts
done
echo "}" >> statusCodes.ts

# Generate Go code
echo "package main" > statusCodes.go
echo "" >> statusCodes.go
echo "type StatusCode int" >> statusCodes.go
echo "" >> statusCodes.go
echo "const (" >> statusCodes.go
for name in "${!status_codes[@]}"; do
    echo "    $name StatusCode = ${status_codes[$name]}" >> statusCodes.go
done
echo ")" >> statusCodes.go