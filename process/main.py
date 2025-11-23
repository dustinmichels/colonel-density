import os

import pandas as pd
from geopy.exc import GeocoderTimedOut
from geopy.extra.rate_limiter import RateLimiter
from geopy.geocoders import Nominatim
from rich.console import Console

console = Console()

INPUT_FILE = "data/locations.csv"
OUTPUT_FILE = "data/locations_updated.csv"
SAVE_EVERY = 5  # save every N successful geocodes
SEARCH_AGAIN = (
    True  # If True, re-search previously searched entries; If False, skip them
)

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

# Add searched column if it doesn't exist
if "searched" not in df.columns:
    df["searched"] = False

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

geolocator = Nominatim(user_agent="kfc_geocoder", timeout=5)

# Use RateLimiter for automatic rate limiting and retries
geocode = RateLimiter(
    geolocator.geocode,
    min_delay_seconds=1,  # Nominatim requires 1 second minimum
    max_retries=3,
    error_wait_seconds=2,
)


def save_progress(df, filename):
    df.drop(columns=["full_address"]).to_csv(filename, index=False)
    console.print(f"[green]Progress saved to[/green] {filename}")


def get_remaining(df):
    if SEARCH_AGAIN:
        return df[(df["latitude"].isna() | df["longitude"].isna())]
    else:
        return df[
            (df["latitude"].isna() | df["longitude"].isna()) & (df["searched"] == False)
        ]


def get_remaining_count(df):
    return len(get_remaining(df))


# ---------------------------
# MAIN LOOP
# ---------------------------

# Filter based on SEARCH_AGAIN setting
missing = get_remaining(df)

console.print(
    f"[bold green]Starting geocoding:[/bold green] {len(missing)} remaining / {len(df)} total"
)

save_counter = 0

for idx, row in missing.iterrows():
    console.print(f"[blue]Geocoding: {row['full_address']}[/blue]")

    try:
        location = geocode(row["full_address"])
    except GeocoderTimedOut:
        console.print("[red]  Error: Geocoder timed out[/red]")
        location = None
    except Exception as e:
        console.print(f"[red]  Error: {e}[/red]")
        location = None

    # Mark as searched regardless of result
    df.at[idx, "searched"] = True

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
                f"[cyan]---------- Progress saved ({get_remaining_count(df)} remaining / {len(df)}) ----------[/cyan]"
            )
            save_counter = 0

    else:
        console.print("[yellow]  → Not found (marked as searched)[/yellow]")

# Final save
save_progress(df, OUTPUT_FILE)
console.print(f"[bold green]Saved final file:[/bold green] {OUTPUT_FILE}")

# final remaining count
console.print(
    f"[bold yellow]Done.[/bold yellow] {get_remaining_count(df)} entries remain un-geocoded out of {len(df)}."
)
