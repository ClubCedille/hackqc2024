# soumission-aide.geojson

## GeoJSON avec la liste des soumissions d'aide reçues dans la province

- `type` : Le type de l'objet GeoJSON, Le type "FeatureCollection" est utilisé.

### Champ `Properties`

- `_id` : Un identifiant unique pour l'enregistrement de la soumission d'aide.
- `source_externe_linked` : Un lien vers la source externe de l'incident pour lequel l'aide est proposée.
- `categorie_catastrophe` : La catégorie de la catastrophe pour laquelle l'aide est proposée.
- `how_to_help` : Instructions ou informations sur comment fournir de l'aide.
- `how_to_use_help` : Instructions ou informations sur comment utiliser l'aide proposée.
- `need_help` : Un booléen indiquant si de l'aide est nécessaire.
- `Type` : Le type de l'objet cartographique.
- `category` : La catégorie de l'aide ou du support proposé (par exemple, "Hébergement").
- `date` : La date associée à l'objet de la carte, au format ISO 8601 (par exemple, "0001-01-01T00:00:00Z").
- `description` : Une description de l'aide ou du support proposé.
- `name` : Le nom ou le titre de l'objet de la carte.

#### Champs `geometry`

- `coordinates` : Un tableau représentant les coordonnées géographiques (longitude, latitude).
- `type` : Le type géométrique, typiquement "Point" pour représenter un emplacement unique.
