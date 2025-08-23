
#!/usr/bin/env bash

# 🧠 Metadata
readonly SCRIPT_NAME="$(basename "$0")"
readonly TIMESTAMP="$(date +%Y-%m-%d_%H-%M-%S)"
LOG_FILE="/tmp/${SCRIPT_NAME}_${TIMESTAMP}.log"

# 🧹 Cleanup on exit
cleanup() {
  echo "Cleaning up..." >> "$LOG_FILE"
}
trap cleanup EXIT

# 📦 Declare an associative array
declare -A status_codes=(
  [200]="OK"
  [404]="Not Found"
  [500]="Server Error"
)

# 📁 Create temp directory
TMP_DIR=$(mktemp -d)
echo "Working in $TMP_DIR" >> "$LOG_FILE"

# 📄 Read and process a file line-by-line
input_file="input.txt"
if [[ ! -f "$input_file" ]]; then
  echo "Error: $input_file not found!" >&2
  exit 1
fi

process_line() {
  local line="$1"
  if [[ "$line" =~ ^[0-9]{3} ]]; then
    code="${line:0:3}"
    message="${status_codes[$code]:-Unknown}"
    echo "[$code] $message - ${line:4}" >> "$LOG_FILE"
  elif [[ "$line" == *"ERROR"* ]]; then
    echo "⚠️ Error detected: $line" >> "$LOG_FILE"
  else
    echo "Info: $line" >> "$LOG_FILE"
  fi
}

while IFS= read -r line; do
  process_line "$line"
done < "${input_file}"

# 🔁 Loop through files
for file in "$TMP_DIR"/*; do
  [[ -e "$file" ]] || continue
  echo "Found file: $file" >> "$LOG_FILE"
done

# 🧪 Run a command and capture output
output=$(curl -s https://example.com || echo "curl failed")
echo "Curl output: ${output:0:50}..." >> "$LOG_FILE"

# Conditional logic
if [[ "$output" == *"Example Domain"* ]]; then
  echo "✅ Site is reachable" >> "$LOG_FILE"
else
  echo "❌ Site unreachable" >> "$LOG_FILE"
fi

# 🧾 Final summary
echo "Log saved to $LOG_FILE"

