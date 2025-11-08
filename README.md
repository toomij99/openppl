
# KFXE Study & Training Bundle (Updated Nov 8, 2025)

Includes:
- PPL_IFR_CPL_Schedule_UPDATED_Nov8_2025.ics  → Import to Calendar
- Merged_Training_Reminders_Updated_Nov8_2025.csv  → Import to Reminders (Shortcuts-based script)
- PPL_Written_Microplan_Reminders_Nov7-21_2025.csv → PPL written micro-plan
- Flight_Training_PPL_IFR_CPL_Tasks_Nov2025-Feb2026.csv → Full master task list
- import_reminders_shortcuts.sh  → CLI importer using Shortcuts (supports real Reminders tags)

## Calendar Import
- Double-click the .ics and add to your training calendar.
- Timezones handled: Asia/Jerusalem (pre-travel) → America/New_York (post-arrival).
- Shabbat blocks included (Fri 16:00 → Sat 18:00 local).

## Reminders Import with Tags (Recommended)
1) Create a Shortcut on your Mac named: **Import Reminder Row**
   - Receive: *Any*
   - Actions:
     1. **Get Dictionary from Input**
     2. **Get Dictionary Value** “list” → List Name
     3. **Get Dictionary Value** “title” → Title
     4. **Get Dictionary Value** “notes” → Notes
     5. **Get Dictionary Value** “due” → Due ISO (Text like `2025-11-09 09:00`)
     6. **Get Dictionary Value** “tags” → Tags (List)
     7. **Add New Reminder** → Title=Title, List=List Name, Notes=Notes, Due Date=Due ISO, Tags=Tags

   - In Shortcuts → Settings → Advanced → enable **Allow Running Scripts**.

2) Run the importer:
```bash
chmod +x import_reminders_shortcuts.sh
./import_reminders_shortcuts.sh -l "Flight Training 2025–2026" -t "09:00" Merged_Training_Reminders_Updated_Nov8_2025.csv
```

You can re-run it for the other CSVs too.

## Notes
- Tags are created as true Reminders tags via Shortcuts.
- Adjust `-t` to change the default due time.
- If you prefer separate lists (PPL / IFR / CPL), run the script multiple times with different `-l` names.
