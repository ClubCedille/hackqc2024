# soumission-aide.json

## JSON avec la liste des soumissions d'aide reçues dans la province

- `_id` : Un identifiant unique pour l'enregistrement de la soumission d'aide.
- `source_externe_linked` : Un lien vers la source externe de l'incident pour lequel l'aide est proposée.
- `categorie_catastrophe` : La catégorie de la catastrophe pour laquelle l'aide est proposée.
- `how_to_help` : Instructions ou informations sur comment fournir de l'aide.
- `how_to_use_help` : Instructions ou informations sur comment utiliser l'aide proposée.
- `need_help` : Un booléen indiquant si de l'aide est nécessaire.

### Champ `map_object`

- `name` : Le titre de l'aide proposé sur la carte.
- `description` : Une description de l'aide proposé.
- `Type` : Le type de l'objet cartographique.
- `category` : La catégorie de l'aide ou du support proposé (par exemple, "Hébergement").
- `date` : La date associée à l'objet de la carte, au format ISO 8601 (par exemple, "0001-01-01T00:00:00Z").

#### Champs `geometry`

- `coordinates` : Les coordonnées géographiques de la soumission d'aide (longitude, latitude).
- `type` : Le type géométrique. Le type "Point" est utilisé pour représenter un emplacement unique.
