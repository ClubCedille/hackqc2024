# soumission-aide.json

## JSON avec la liste des soumissions d'aide reçues dans la province

- `_id` : Un identifiant unique pour l'enregistrement, généralement un UUID.
- `contact_infos` : Informations de contact, typiquement un numéro de téléphone.
- `event_id` : Un identifiant pour l'événement.
- `how_to_help` : Instructions ou informations sur comment fournir de l'aide.
- `how_to_use_help` : Instructions ou informations sur comment utiliser l'aide proposée.
- `need_help` : Un booléen indiquant si de l'aide est nécessaire.

### Champ `map_object`

- `Id` : Un identifiant pour l'objet de la carte.
- `Type` : Le type de l'objet cartographique.
- `account_id` : Un identifiant unique, typiquement pour un compte, associé à l'objet de la carte.
- `category` : La catégorie de l'aide ou du support proposé (par exemple, "Hébergement").
- `date` : La date associée à l'objet de la carte, au format ISO 8601 (par exemple, "0001-01-01T00:00:00Z").
- `description` : Une description de l'aide ou du support proposé.
- `name` : Le nom ou le titre de l'objet de la carte.

#### Champs `geometry`

- `coordinates` : Un tableau représentant les coordonnées géographiques (longitude, latitude).
- `type` : Le type géométrique, typiquement "Point" pour représenter un emplacement unique.
