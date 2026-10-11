#!/usr/bin/env bash
set -euo pipefail

DEVICE_SERIAL="${1:-0123456789ABCDEF}"
API_URL="${2:-http://127.0.0.1:8080}"
API_KEY="${JERVIS_API_AUTH_KEY:-op-secret-key-98765}"

echo "=== Alcatel 5022D SIM & Carrier Monitor ==="
echo "Target Device Serial: $DEVICE_SERIAL"
echo "Jervis Gateway URL:   $API_URL"
echo ""

check_sim() {
    adb -s "$DEVICE_SERIAL" shell getprop gsm.sim.state | tr -d '\r'
}

STATE=$(check_sim)
echo "Current SIM State: $STATE"

if [[ "$STATE" == *"ABSENT"* ]]; then
    echo ""
    echo "[!] Physical SIM card is NOT detected."
    echo "    Please power off the phone, remove back cover/battery,"
    echo "    insert an active GSM SIM card into SIM Slot 1, and power back on."
    echo "    Waiting for SIM state to become READY..."
    echo ""

    while true; do
        STATE=$(check_sim)
        if [[ "$STATE" == *"READY"* ]] || [[ "$STATE" == *"LOADED"* ]]; then
            echo "[+] SIM detected! New state: $STATE"
            break
        fi
        sleep 3
    done
fi

OPERATOR=$(adb -s "$DEVICE_SERIAL" shell getprop gsm.operator.alpha | tr -d '\r')
echo "[+] Detected Carrier / Operator: $OPERATOR"

read -p "Enter destination phone number to test outbound SMS (e.g. +1234567890): " RECIPIENT

CMD_ID="cmd-real-sms-$(date +%s)"
PAYLOAD=$(cat <<EOF
{
  "schema": "phone-agent.v1",
  "id": "$CMD_ID",
  "device": "alcatel-5022d",
  "timestamp": $(date +%s000),
  "type": "command",
  "payload": {
    "command_id": "test-sms-$(date +%s)",
    "action": "send_sms",
    "recipient": "$RECIPIENT",
    "text": "Jervis -> 5022D Hardware Link Verified."
  }
}
EOF
)

echo "[*] Dispatching send_sms command to Jervis Gateway..."
curl -s -X POST "$API_URL/v1/device/commands" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD"

echo ""
echo "[*] Waiting for execution result from physical device..."
for i in {1..15}; do
    RESULT=$(curl -s -H "Authorization: Bearer $API_KEY" "$API_URL/v1/device/status?device=alcatel-5022d")
    echo "Status check $i: $RESULT"
    sleep 2
done
