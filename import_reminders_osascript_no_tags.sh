#!/usr/bin/env bash
# import_reminders_osascript_no_tags.sh
# Usage:
#   ./import_reminders_osascript_no_tags.sh -l "Flight Training 2025–2026" -t "09:00" Reminders_ALL_MERGED_WeekNumbers_Nov8_2025.csv
#
# Notes:
#   - Creates or uses the given Reminders list.
#   - Sets due date with a uniform time (HH:MM). Ignores "Tag" column.
#   - Uses AppleScript via `osascript` (no Shortcuts, no tags).

set -euo pipefail

LIST_NAME=""
DUE_TIME="09:00"

usage() {
  echo "Usage: $0 -l \"List Name\" [-t HH:MM] file.csv" >&2
  exit 64
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
[[ -z "$LIST_NAME" || -z "$CSV_PATH" ]] && usage
[[ -f "$CSV_PATH" ]] || { echo "CSV not found: $CSV_PATH" >&2; exit 66; }
[[ "$DUE_TIME" =~ ^([0-1][0-9]|2[0-3]):[0-5][0-9]$ ]] || { echo "Bad -t time (HH:MM)"; exit 65; }

/usr/bin/python3 - "$CSV_PATH" "$LIST_NAME" "$DUE_TIME" <<'PY'
import csv, sys, subprocess, datetime

csv_path, list_name, due_time = sys.argv[1], sys.argv[2], sys.argv[3]
hh, mm = map(int, due_time.split(":"))
seconds_from_midnight = hh*3600 + mm*60

APPLESCRIPT = r"""
on run argv
  set listName to item 1 of argv
  set secondsFromMidnight to (item 2 of argv) as integer
  set y to (item 3 of argv) as integer
  set m to (item 4 of argv) as integer
  set d to (item 5 of argv) as integer
  set theTitle to item 6 of argv
  set theNotes to item 7 of argv
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
    make new reminder at end of reminders of theList with properties {name:theTitle, body:theNotes, due date:theDate}
  end tell
end run
"""

def parse_date(s):
  for fmt in ("%Y-%m-%d","%m/%d/%Y"):
    try:
      return datetime.datetime.strptime(s, fmt).date()
    except Exception:
      pass
  return None

with open(csv_path, newline='', encoding='utf-8') as f:
  r = csv.reader(f)
  header = next(r, [])
  for row in r:
    if not row or len(row) < 2:
      continue
    due_s = (row[0] or "").strip()
    title = (row[1] or "").strip()
    notes = (row[3] or "").strip() if len(row)>3 else ""
    if not due_s or not title:
      continue
    dt = parse_date(due_s)
    if not dt:
      print(f"Skipping row with unrecognized date: {due_s}", file=sys.stderr)
      continue
    args = [
      "osascript",
      "-e", APPLESCRIPT,
      "--",
      list_name,
      str(seconds_from_midnight),
      str(dt.year), str(dt.month), str(dt.day),
      title,
      notes,
    ]
    subprocess.run(args, check=True)
print("Import complete.")
PY

echo "✅ Imported (no tags) into list: $LIST_NAME at $DUE_TIME"
