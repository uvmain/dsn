# Simple API Test
$headers = @{'Content-Type' = 'application/json'}

# Test 1: Register admin user
Write-Host "Testing user registration..." -ForegroundColor Yellow
$body = @{
    username = "admin"
    email = "admin@example.com"  
    password = "admin123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/register" -Method POST -Body $body -Headers $headers -SessionVariable session
    Write-Host "Success! Registered user: $($response.username), Admin: $($response.is_admin)" -ForegroundColor Green
} catch {
    Write-Host "Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 2: Create a note  
Write-Host "Testing note creation..." -ForegroundColor Yellow
$noteBody = @{
    title = "Test Note"
    content = "This is my first note"
    color = "#ffeb3b"
} | ConvertTo-Json

try {
    $noteResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/notes" -Method POST -Body $noteBody -Headers $headers -WebSession $session
    Write-Host "Success! Created note: $($noteResponse.title) with ID: $($noteResponse.id)" -ForegroundColor Green
} catch {
    Write-Host "Failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "Test completed!" -ForegroundColor Cyan
