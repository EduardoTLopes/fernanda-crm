// Replace YOUR_API_KEY with your actual API key
const apiKey = process.env.GOOGLE_MAPS_API_KEY;
const location = { lat: -23.570552, lng: -46.635314 };
const radius = 3000; // Search radius in meters
const type = "bakery";
const fields = "formatted_phone_number"; // Include phone numbers in the response

const generalUrl = `https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=${location.lat},${location.lng}&radius=${radius}&type=${type}&fields=${fields}&key=${apiKey}`;

// Make a request to the Places API
let places = [];
async function fetchPlaces() {
  try {
    const response = await fetch(generalUrl);
    const data = await response.json();
    const results = data.results;
    for (const result of results) {
      let place = {
        id: result.place_id,
        name: result.name,
        phone: "",
      };
      const specificUrl = `https://maps.googleapis.com/maps/api/place/details/json?place_id=${place.id}&fields=${fields}&key=${apiKey}`;
      try {
        const phoneResponse = await fetch(specificUrl);
        const phoneData = await phoneResponse.json();
        place.phone = phoneData.result.formatted_phone_number;
        places.push(place);
      } catch (error) {
        console.error("Error fetching phone number:", error);
      }
    }
  } catch (error) {
    console.error("Error fetching data:", error);
  } finally {
    places.forEach((place) =>
      console.log({ name: place.name, phone: place.phone })
    );
  }
}

fetchPlaces();
