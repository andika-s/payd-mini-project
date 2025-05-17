# Payd Mini Project: Daily Worker Roster Management System

## Project Overview

A scheduling system for daily workers with employee and admin capabilities

### Employee Features
- View assigned shifts (approved requests with status tracking)
- Browse available shifts (unassigned shifts with role filters)
- Submit shift requests (creates pending status for approval)
- Track request status (Pending/Approved/Rejected indicators)

### Admin Features
- Shift schedule management (Create/update/delete shifts)
- Request processing (Approve/reject worker requests)
- Roster assignments (Manage worker-shift allocations)

## Setup Using Docker Compose

1. Clone repository:
git clone <repository-url>
cd payd-mini-project

2. Start services:
docker-compose up --build

Access endpoints:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## Documented Assumptions

1. Worker Identification
   - Persistent unique worker IDs
   - ID-based authentication system
   - No password/PIN requirements

2. Shift Management
   - UTC timezone storage
   - Max 5 shifts/week per worker
   - No overlapping shifts allowed

3. Request Handling
   - 24h request expiration
   - Single active request per shift
   - First-come-first-served basis

## API Usage Examples

Get available shifts:
curl http://localhost:8080/api/v1/shifts

Request a shift:
curl -X POST -H "Content-Type: application/json" -d '{"shift_id":1535293522186729124,"worker_id":1000}' http://localhost:8080/api/v1/shift/request

## UI Walkthrough Guide

1. View Assigned Shifts
   - Navigate to "My Shifts" tab
   - Approved shifts shown with status badges
   - Color indicators: Green = Approved, Yellow = Pending

2. Request Available Shifts
   - Select "Available Shifts" tab
   - Click "Request Shift" on desired shift card
   - Button updates to "Requested" during processing

3. Track Requests
   - Pending requests visible in "My Shifts"
   - Automatic refresh on status changes
   - System notifications for approvals/rejections