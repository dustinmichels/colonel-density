# run scraper
cd scraper
go run .
cd ..

# copy output
cp scraper/out/locations.csv process/data/locations.csv

cd process
uv run main.py

