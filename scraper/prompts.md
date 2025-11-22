# Prompts

## Locations

Next, iterate over the locations.

locations, err := getLocationsOnStatePage(url, stateCode)

Visit each location URL and extract the addresses of the KFC's there.

This getAddressesFromPlaceshould go in another file.

The address will either be in HTML like this:

```html
<div class="Core-address">
  <span
    class="coordinates"
    itemprop="geo"
    itemscope=""
    itemtype="http://schema.org/GeoCoordinates"
    ><meta itemprop="latitude" content="38.71639645751588" /><meta
      itemprop="longitude"
      content="-121.39278798834323"
  /></span>
  <address
    class="c-address"
    id="address"
    itemscope=""
    itemtype="http://schema.org/PostalAddress"
    itemprop="address"
    data-country="US"
  >
    <meta itemprop="addressLocality" content="Antelope" /><meta
      itemprop="streetAddress"
      content="8101 Watt Avenue"
    />
    <div class="c-AddressRow">
      <span class="c-address-street-1">8101 Watt Avenue</span>
    </div>
    <div class="c-AddressRow">
      <span class="c-address-city">Antelope</span><yxt-comma>,</yxt-comma>
      <abbr
        title="California"
        aria-label="California"
        class="c-address-state"
        itemprop="addressRegion"
        >CA</abbr
      >
      <span class="c-address-postal-code" itemprop="postalCode">95843</span>
    </div>
    <div class="c-AddressRow">
      <abbr
        title="United States"
        aria-label="United States"
        class="c-address-country-name c-address-country-us"
        itemprop="addressCountry"
        >US</abbr
      >
    </div>
  </address>
</div>
```
