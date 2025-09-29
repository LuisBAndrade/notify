#!/bin/bash

USER_ID="123"
API="http://localhost:8080/api"

echo "👉 Sending notification..."
SEND_RESPONSE=$(curl -s -X POST $API/notification/send \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'","payload":{"text":"Hello safe parse"}}')

echo "Response: $SEND_RESPONSE"

# Try extracting message_id (but don’t exit on failure)
MESSAGE_ID=$(echo "$SEND_RESPONSE" | grep -o '"message_id":"[^"]*"' | cut -d':' -f2 | tr -d '"')

if [ -z "$MESSAGE_ID" ]; then
  echo "⚠️ Could not parse message_id, continuing..."
else
  echo "Extracted message_id: $MESSAGE_ID"
  
  echo "👉 Acknowledging..."
  curl -s -X POST $API/notification/acknowledge \
    -H "Content-Type: application/json" \
    -d '{"user_id":"'"$USER_ID"'","message_ids":["'"$MESSAGE_ID"'"]}'
  echo ""
fi