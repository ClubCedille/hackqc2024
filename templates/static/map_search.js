const { autocomplete } = window['@algolia/autocomplete-js'];

var temporaryMarker;

/**
 * @param {string} address partial address
 * @returns {Array<string>} List of addresses
 */
async function getAdresses(address) {
    const res = await fetch(`https://geoegl.msp.gouv.qc.ca/apis/icherche/geocode?type=adresses&limit=5&geometry=true&q=${address}`);
    const json = await res.json();
    return json.features.map(feature => {
        return feature
    });
}

async function getMuninipalites(address) {
    const res = await fetch(`https://geoegl.msp.gouv.qc.ca/apis/icherche/geocode?type=municipalites&limit=5&geometry=true&q=${address}`);
    const json = await res.json();
    return json.features.map(feature => {
        return feature;
    });
}

function multiPolygonToCoords(feature){
    let coords = feature.geometry.coordinates;
    while (coords.length != 2) {
        coords = coords[0];
    }
    return coords;
}

function dropdownOnSelect(event) {
    let geojson = L.geoJSON(event.item);
    flyToItem(geojson);
    addTemporaryMarker(geojson);
}

function flyToItem(geojson) {
    let bounds = geojson.getBounds();
    map.flyToBounds(bounds, {maxZoom: 15});
}

function addTemporaryMarker(geojson) {
    let title = geojson.getLayers()[0].feature.properties.nom;
    if(temporaryMarker){
        map.removeLayer(temporaryMarker);
    }
    temporaryMarker = L.marker(geojson.getBounds().getCenter() , 
    {
        draggable: false,
        title: title,
        alt: "Centre de la carte",

    }).bindTooltip(title, 
    {
        permanent: true, 
        direction: 'right'
    });
    temporaryMarker.addTo(map);
}

autocomplete({
    container: '#autocomplete',
    placeholder: 'Adresse...',
    getSources({ query }) {
        return [
            {
                sourceId: 'addresse',
                async getItems() {
                    let addresses = await getAdresses(query);
                    return addresses;
                },
                templates: {
                    item({ item }) {
                        return `${item.properties.nom}`;
                    },
                },
                onSelect: function (event) {
                    dropdownOnSelect(event)
                },
            },
            {
                sourceId: 'municipalite',
                async getItems() {
                    let muninipalites = await getMuninipalites(query);
                    return muninipalites;
                },
                templates: {
                    item({ item }) {
                        return `${item.properties.nom}`;
                    },
                },
                onSelect: function (event) {
                    dropdownOnSelect(event)
                },
            }
        ];
    },
});
