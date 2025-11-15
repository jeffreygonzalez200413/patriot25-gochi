import datetime
from typing import List
from googleapiclient.discovery import build
from google.oauth2.credentials import Credentials

SCOPES = ["https://www.googleapis.com/auth/calendar.readonly"]
TOKEN_FILE = "./token.json"


def get_upcoming_events(max_results: int = 5) -> List[str]:
    """
    Returns a list of human-readable event strings like:
    'Today 15:00-16:00: Project meeting'
    """
    creds = Credentials.from_authorized_user_file(TOKEN_FILE, SCOPES)
    service = build("calendar", "v3", credentials=creds)

    now = datetime.datetime.utcnow().isoformat() + "Z"  # 'Z' = UTC
    events_result = (
        service.events()
        .list(
            calendarId="primary",
            timeMin=now,
            maxResults=max_results,
            singleEvents=True,
            orderBy="startTime",
        )
        .execute()
    )
    events = events_result.get("items", [])

    summaries: List[str] = []
    for event in events:
        start = event.get("start", {})
        end = event.get("end", {})
        summary = event.get("summary", "(no title)")

        start_str = start.get("dateTime") or start.get("date")
        end_str = end.get("dateTime") or end.get("date")

        # simplify timestamps for prompt (rough local-ish view)
        if start_str and "T" in start_str:
            start_str = start_str.replace("T", " ").split("+")[0]
        if end_str and "T" in end_str:
            end_str = end_str.replace("T", " ").split("+")[0]

        summaries.append(f"{start_str} ~ {end_str}: {summary}")

    return summaries
