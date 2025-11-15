from calendar_client import get_upcoming_events

if __name__ == "__main__":
    try:
        events = get_upcoming_events(5)
        print(f"Got {len(events)} events")
        for ev in events:
            print("-", ev)
    except Exception as e:
        import traceback

        traceback.print_exc()
        print("Error fetching events:", e)
