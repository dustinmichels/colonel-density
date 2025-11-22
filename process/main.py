import time

import pandas as pd
from geopy.exc import GeocoderTimedOut, GeocoderUnavailable
from geopy.geocoders import Nominatim

df = pd.read_csv("data/locations.csv")

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
TIMEOUT = 10  # seconds
RETRIES = 5  # attempts per address
DELAY = 1.2  # seconds (required by Nominatim)


def safe_geocode(address):
    for attempt in range(1, RETRIES + 1):
        try:
            return geolocator.geocode(address, timeout=TIMEOUT)
        except (GeocoderTimedOut, GeocoderUnavailable):
            print(f"  Timeout (attempt {attempt}/{RETRIES}) → retrying...")
            time.sleep(2)  # wait a bit longer before retrying
        except Exception as e:
            print(f"  Unexpected error: {e}")
            return None
    return None  # give up


missing = df[df["latitude"].isna() | df["longitude"].isna()]

for idx, row in missing.iterrows():
    addr = row["full_address"]
    print(f"Geocoding: {addr}")

    location = safe_geocode(addr)
    time.sleep(DELAY)

    if location:
        df.at[idx, "latitude"] = location.latitude
        df.at[idx, "longitude"] = location.longitude
        print(f"  → Found: {location.latitude}, {location.longitude}")
    else:
        print("  → Failed after retries")

# df = df.drop(columns=["full_address"])
df.to_csv("data/out/locations_geocoded.csv", index=False)
print("Saved as data/out/locations_geocoded.csv")
