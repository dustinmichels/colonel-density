import os
import time

import pandas as pd
from geopy.exc import GeocoderTimedOut, GeocoderUnavailable
from geopy.geocoders import Nominatim

INPUT_FILE = "data/locations.csv"
OUTPUT_FILE = "data/locations_updated.csv"

SAVE_EVERY = 5  # save every N successful geocodes

# ---------------------------
# Load existing progress first
# ---------------------------
if os.path.exists(OUTPUT_FILE):
    print(f"Resuming from existing file: {OUTPUT_FILE}")
    df = pd.read_csv(OUTPUT_FILE)
else:
    df = pd.read_csv(INPUT_FILE)

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
            return geolocator.geocode(address, timeout=TIMEOUT)
        except (GeocoderTimedOut, GeocoderUnavailable):
            print(f"  Timeout (attempt {attempt}/{RETRIES}) — retrying...")
            time.sleep(2)
        except Exception as e:
            print(f"  Unexpected error: {e}")
            return None
    return None


def save_progress(df, filename):
    df.drop(columns=["full_address"]).to_csv(filename, index=False)
    print(f"Progress saved to {filename}")


def get_remaining_count(df):
    return df["latitude"].isna().sum() + df["longitude"].isna().sum()


# ---------------------------
# MAIN LOOP
# ---------------------------

missing = df[df["latitude"].isna() | df["longitude"].isna()]
total_missing_start = len(missing)
total_rows = len(df)

# print initial remaining count
print(f"Starting geocoding: {total_missing_start} remaining / {total_rows} total")

save_counter = 0

for idx, row in missing.iterrows():
    print(f"Geocoding: {row['full_address']}")

    location = safe_geocode(row["full_address"])
    time.sleep(DELAY)

    if location:
        df.at[idx, "latitude"] = location.latitude
        df.at[idx, "longitude"] = location.longitude
        print(f"  → Found: {location.latitude}, {location.longitude}")

        save_counter += 1

        if save_counter >= SAVE_EVERY:
            # save file
            save_progress(df, OUTPUT_FILE)

            # print remaining count at save time
            print(
                f"---------- Progress saved ({get_remaining_count(df)} remaining / {total_rows}) ----------"
            )

            save_counter = 0

    else:
        print("  → Failed after retries")

# Final save
save_progress(df, OUTPUT_FILE)
print(f"Saved final file: {OUTPUT_FILE}")

# final remaining count
print(
    f"Done. {get_remaining_count(df)} entries remain un-geocoded out of {total_rows}."
)
