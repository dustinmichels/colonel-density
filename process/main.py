import os
import time

import pandas as pd
from geopy.exc import GeocoderTimedOut, GeocoderUnavailable
from geopy.geocoders import Nominatim
from rich.console import Console

console = Console()

INPUT_FILE = "data/locations.csv"
OUTPUT_FILE = "data/locations_updated.csv"
SAVE_EVERY = 5  # save every N successful geocodes

# ---------------------------
# Load existing progress first
# ---------------------------
if os.path.exists(OUTPUT_FILE):
    console.print(f"[cyan]Resuming from existing file:[/cyan] {OUTPUT_FILE}")
    df = pd.read_csv(OUTPUT_FILE)
else:
    df = pd.read_csv(INPUT_FILE)

# sort df by state, then city
df = df.sort_values(by=["state", "city"]).reset_index(drop=True)

# Build full address if missing
if "full_address" not in df.columns:
    df["full_address"] = (
        df["address"].astype(str)
        + ", "
        + df["city"].astype(str)
        + ", "
        + df["state"].astype(str)
        + " "
        + df["zip_code"].astype(str)
        + ", "
        + df["country"].astype(str)
    )

geolocator = Nominatim(user_agent="kfc_geocoder")
TIMEOUT = 10
RETRIES = 5
DELAY = 1.2


def safe_geocode(address):
    for attempt in range(1, RETRIES + 1):
        try:
            geo = geolocator.geocode(address, timeout=TIMEOUT)
            if geo:
                return geo
            else:
                console.print("[yellow]  No result found for address.[/yellow]")
                return None
        except (GeocoderTimedOut, GeocoderUnavailable):
            console.print(
                f"[red]  Timeout (attempt {attempt}/{RETRIES}) — retrying...[/red]"
            )
            time.sleep(2)
        except Exception as e:
            console.print(f"[red]  Unexpected error: {e}[/red]")
            return None
    return None


def save_progress(df, filename):
    df.drop(columns=["full_address"]).to_csv(filename, index=False)
    console.print(f"[green]Progress saved to[/green] {filename}")


def get_remaining_count(df):
    return len(df[df["latitude"].isna() | df["longitude"].isna()])


# ---------------------------
# MAIN LOOP
# ---------------------------

missing = df[df["latitude"].isna() | df["longitude"].isna()]
total_missing_start = len(missing)
total_rows = len(df)

console.print(
    f"[bold green]Starting geocoding:[/bold green] {total_missing_start} remaining / {total_rows} total"
)

save_counter = 0

for idx, row in missing.iterrows():
    console.print(f"[blue]Geocoding: {row['full_address']}[/blue]")

    location = safe_geocode(row["full_address"])
    time.sleep(DELAY)

    if location:
        df.at[idx, "latitude"] = location.latitude
        df.at[idx, "longitude"] = location.longitude

        console.print(
            f"[green]  → Found:[/green] {location.latitude}, {location.longitude}"
        )

        save_counter += 1

        if save_counter >= SAVE_EVERY:
            save_progress(df, OUTPUT_FILE)
            console.print(
                f"[cyan]---------- Progress saved ({get_remaining_count(df)} remaining / {total_rows}) ----------[/cyan]"
            )
            save_counter = 0

    else:
        console.print("[red]  → Failed after retries[/red]")

# Final save
save_progress(df, OUTPUT_FILE)
console.print(f"[bold green]Saved final file:[/bold green] {OUTPUT_FILE}")

# final remaining count
console.print(
    f"[bold yellow]Done.[/bold yellow] {get_remaining_count(df)} entries remain un-geocoded out of {total_rows}."
)
