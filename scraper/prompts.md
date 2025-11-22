# Prompts

## Locations

Next, iterate over the cities and get locations.

Visit each city URL and extract the addresses of the locations there.

This `getLocationsFromCity` function should go in a new file. Location should include name, address, lat/lon (if available).

The City.DataCount field tells you how many locations are in a city. For each city scraped, not how many locations could be successfully parsed (eg, "successfully got 2/2 locations").

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

Or like this:

```html
<li class="Directory-listTeaser">
  <article class="Teaser Teaser--ace Teaser--directory">
    <h2 class="Teaser-title" aria-level="2">
      <a
        class="Teaser-titleLink"
        href="../ca/antioch/2410-mahogany-way"
        data-ya-track="businessname"
        ><span class="LocationName LocationName--directory"
          ><span class="LocationName-brand">KFC</span>
          <span class="LocationName-geo">2410 Mahogany Way</span></span
        ></a
      >
    </h2>
    <div class="Teaser-open">
      <span
        class="c-hours-today js-hours-today"
        data-days='[{"day":"MONDAY","intervals":[],"isClosed":true},{"day":"TUESDAY","intervals":[],"isClosed":true},{"day":"WEDNESDAY","intervals":[],"isClosed":true},{"day":"THURSDAY","intervals":[],"isClosed":true},{"day":"FRIDAY","intervals":[],"isClosed":true},{"day":"SATURDAY","intervals":[],"isClosed":true},{"day":"SUNDAY","intervals":[],"isClosed":true}]'
        data-utc-offsets='[{"offset":-28800,"start":1762074000},{"offset":-25200,"start":1772964000},{"offset":-28800,"start":1793523600}]'
        ><span class="Hours-status Hours-status--ace Hours-status--closed"
          >Closed</span
        ></span
      >
    </div>
    <div class="Teaser-address">
      <address class="c-address" data-country="US">
        <div class="c-AddressRow">
          <span class="c-address-street-1">2410 Mahogany Way</span>
        </div>
        <div class="c-AddressRow">
          <span class="c-address-city">Antioch</span><yxt-comma>,</yxt-comma>
          <abbr
            title="California"
            aria-label="California"
            class="c-address-state"
            >CA</abbr
          >
          <span class="c-address-postal-code">94509</span>
        </div>
        <div class="c-AddressRow">
          <abbr
            title="United States"
            aria-label="United States"
            class="c-address-country-name c-address-country-us"
            >US</abbr
          >
        </div>
      </address>
    </div>
    <div class="Teaser-phone">
      <div class="Phone Phone--main">
        <div class="Phone-label"><span class="sr-only">phone</span></div>
        <div class="Phone-numberWrapper">
          <div class="Phone-display Phone-display--withLink" id="phone-main">
            (925) 754-3474
          </div>
          <div class="Phone-linkWrapper">
            <a class="Phone-link" href="tel:+19257543474" data-ya-track="phone"
              >(925) 754-3474</a
            >
          </div>
        </div>
      </div>
    </div>
    <div class="Teaser-services">
      <div class="Teaser-servicesLabel">Services</div>
      Carry Out, Delivery, No-Contact Delivery
    </div>
    <div class="Teaser-linksRow">
      <div class="Teaser-links">
        <div class="Teaser-link Teaser-directions">
          <div class="c-get-directions">
            <div class="c-get-directions-button-wrapper">
              <a
                class="c-get-directions-button"
                href="http://maps.google.com/?q=2410+Mahogany+Way%2C+Antioch%2C+CA+94509+US&amp;output=classic"
                target="_blank"
                rel="nofollow noopener noreferrer"
                data-ga-category="Get Directions"
                data-ya-track="directions"
                >Get Directions<span class="sr-only wcag-new-tab-hover">
                  Link Opens in New Tab</span
                ></a
              >
            </div>
          </div>
        </div>
        <div class="Teaser-link Teaser-cta">
          <a data-ya-track="visitpage" href="../ca/antioch/2410-mahogany-way"
            >Visit Store Website</a
          >
        </div>
      </div>
    </div>
  </article>
</li>
```
