#!/usr/bin/env bash
# import_reminders.sh
# Usage:
#   ./import_reminders.sh -l "Flight Training 2025–2026" -t "09:00" PPL_Written_Microplan_Reminders_Nov7-21_2025.csv
#   ./import_reminders.sh -l "Flight Training 2025–2026" Flight_Training_PPL_IFR_CPL_Tasks_Nov2025-Feb2026.csv
#
# Flags:
#   -l  Reminders list name (required)
#   -t  Due time (HH:MM, default 09:00)
#   -h  Help

set -euo pipefail

LIST_NAME=""
DUE_TIME="09:00"

usage() {
  sed -n '1,40p' "$0" | sed 's/^# \{0,1\}//'
  exit 1
}

while getopts ":l:t:h" opt; do
  case "$opt" in
    l) LIST_NAME="$OPTARG" ;;
    t) DUE_TIME="$OPTARG" ;;
    h|*) usage ;;
  esac
done
shift $((OPTIND - 1))

CSV_PATH="${1:-}"
if [[ -z "$LIST_NAME" || -z "$CSV_PATH" ]]; then
  echo "Error: list name (-l) and CSV path are required." >&2
  usage
fi
if [[ ! -f "$CSV_PATH" ]]; then
  echo "Error: CSV not found at '$CSV_PATH'." >&2
  exit 1
fi

# Validate DUE_TIME HH:MM
if ! [[ "$DUE_TIME" =~ ^([0-1][0-9]|2[0-3]):[0-5][0-9]$ ]]; then
  echo "Error: Due time must be HH:MM (24h), got '$DUE_TIME'." >&2
  exit 1
fi

/usr/bin/python3 - <<'PY' "$CSV_PATH" "$LIST_NAME" "$DUE_TIME"
import csv, sys, subprocess, datetime

csv_path, list_name, due_time = sys.argv[1], sys.argv[2], sys.argv[3]
hh, mm = map(int, due_time.split(":"))
seconds_from_midnight = hh*3600 + mm*60

# AppleScript program that will be fed each row via argv (after the --)
APPLESCRIPT = r'''
on run argv
  set listName to item 1 of argv
  set secondsFromMidnight to (item 2 of argv) as integer
  set y to (item 3 of argv) as integer
  set m to (item 4 of argv) as integer
  set d to (item 5 of argv) as integer
  set theTitle to item 6 of argv
  set theNotes to item 7 of argv
  set theTag to item 8 of argv
  tell application "Reminders"
    if not (exists list listName) then
      make new list with properties {name:listName}
    end if
    set theList to list listName
    set theDate to current date
    set year of theDate to y
    set month of theDate to (item m of {January, February, March, April, May, June, July, August, September, October, November, December})
    set day of theDate to d
    set time of theDate to secondsFromMidnight
    set combinedNotes to theNotes
    if (length of theTag) > 0 then
      set combinedNotes to theNotes & return & "#" & theTag
    end if
    make new reminder at end of reminders of theList with properties {name:theTitle, body:combinedNotes, due date:theDate}
  end tell
end run
'''

def is_header(row):
    if not row: return True
    return row[0].strip().lower() in ("due date","due","date")

with open(csv_path, newline='', encoding='utf-8') as f:
    reader = csv.reader(f)
    for i, row in enumerate(reader):
        if not row or len(row) < 2:
            continue
        if i == 0 and is_header(row):
            continue

        # Expect columns: Due Date, Title, Tag, Notes
        # Allow shorter rows; fill missing fields with ""
        due_str = (row[0] or "").strip()
        title   = (row[1] or "").strip()
        tag     = (row[2] or "").strip() if len(row) > 2 else ""
        notes   = (row[3] or "").strip() if len(row) > 3 else ""

        if not due_str or not title:
            continue

        # Parse YYYY-MM-DD
        try:
            y, m, d = map(int, due_str.split("-"))
        except Exception:
            # Try other common formats (MM/DD/YYYY)
            try:
                dt = datetime.datetime.strptime(due_str, "%m/%d/%Y")
                y, m, d = dt.year, dt.month, dt.day
            except Exception:
                print(f"Skipping row with unrecognized date: {due_str}", file=sys.stderr)
                continue

        # Call osascript with arguments after a "--"
        args = [
            "osascript",
            "-e", APPLESCRIPT,
            "--",
            list_name,
            str(seconds_from_midnight),
            str(y), str(m), str(d),
            title,
            notes,
            tag,
        ]
        # Important: let osascript handle Unicode; avoid shell quoting by passing as args list.
        subprocess.run(args, check=True)
print("Import complete.")
PY

echo "✅ Imported reminders from '$CSV_PATH' into list '$LIST_NAME' at due time $DUE_TIME."
echo "Note: On first run, macOS may ask to allow Terminal to access Reminders."

