[![Go Report Card](https://goreportcard.com/badge/github.com/Terag/kadok)](https://goreportcard.com/report/github.com/Terag/kadok)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/Terag/kadok)
[![Release](https://img.shields.io/github/v/release/Terag/kadok.svg?style=flat-square)](https://github.com/Terag/kadok/releases/latest)

# Who is Kadok ?

Kadok is a bot developed for the Discord Guild "Les petits pedestres".
It aims to provide fun and useful functionnalities for the Guild Members.

The code of Kadok is open source, you can use it and modify it as long as you follow
the GNU_v3 licence terms.

## Motivation

You like the French serie Kaamelott, own a Discord Guild and would like to add more
features to it ? Don't think twice, this bot is for you!

## Features

The current available features are:
- Play ping pong
- Quote characters from Kaamelott universe (and even more!)
- Manage features' permission based on Discord roles

# Installation

## prerequisites

Before to think about running Kadok in your environment, you must create an application
on Discord's developer portal: https://discord.com/developers/applications

For Discord documentation on how to create an application and add it to your Guild see: https://discord.com/developers/docs/intro

## As a go package

To install the bot locally you can run:

~~~ sh
  go install github.com/Terag/kadok
~~~

## Configuration

Roles and basic properties are customizable. The roles structure should be set in a way that reflects your Discord's guild roles hierarchy.

### Global

The global configuration of the bot is stored in a single yaml file.
The default path for this file being: _properties.yaml_

This file old basic values to initialize modules. The values are separated between the packages. Each package hold the responsibility of defining the properties format.

Example:
```yaml
security:
  roles: "security/roles.json"
characters:
  folder: "characters/resources"
```

### Packages

#### Security

The security package is one of the most important package to configure since it will ensure that the bot works well with your guild. The name of the roles that Kadok MUST be the same from what your guild uses. However, you can create non discord roles in order to make full use of the permissions inheritance.

Available values for Permissions:

 - GetHelp
 - CallCharacter
 - GetCharacterList

Example:

```json
{
    "roles": [
        {
            "name": "member",
            "parent": null,
            "permissions": ["CallCharacter", "GetHelp", "GetCharacterList"]
        },
        {
            "name": "moderator",
            "parent": "member",
            "permissions": []
        },
        {
            "name": "administrator",
            "parent": "moderator",
            "permissions": []
        },
        {
            "name": "Mod - La Table Ronde",
            "parent": "moderator",
            "permissions": []
        },
        {
            "name": "Mod - Bartender",
            "parent": "moderator",
            "permissions": []
        },
        {
            "name": "Mod - Croupier",
            "parent": "moderator",
            "permissions": []
        },
        {
           "name": "Admin - Conseil de Guerre",
           "parent": "administrator",
           "permissions": []
        },
        {
            "name": "Pédèstres seniors",
            "parent": "member",
            "permissions": []
        },
        {
            "name": "Semi Croustillants",
            "parent": "member",
            "permissions": []
        },
        {
            "name": "Indépendentistes Gallois",
            "parent": "member",
            "permissions": []
        },
        {
            "name": "Team Chateau",
            "parent": "member",
            "permissions": []
        }
    ]
}
```

#### Characters

The Characters package bundle all the features related to the quoting of characters.

The properties are:

 - **folder:** The path of the folder from which the characters are generated.

A character is generated from a json file and will take the name of the file. For example, arthur.json becomes arthur for the bot.

Example for arthur.json

```json
{
    "name": "arthur",
    "sentences": [
        "C’est vrai ce qu’on dit, vous êtes le fils d’un démon et d’une pucelle ? Vous avez plus pris de la pucelle.",

        "Ça vous fait pas mal à la tête de glandouiller vingt-quatre heures sur vingt-quatre ?",

        "Y a pas à dire, dès qu'il y a du dessert, le repas est tout de suite plus chaleureux !"
    ]
}
```

## Run

### Simple


You can run it by calling it from your shell:

~~~ shell script
kadok.exe -t <bot-token> (-p <properties-file>)
~~~

Available options:

 - -t \<bot-token\> : **required**, bot's token to access discord's api
 - -p \<properties-file\> : optionnal, bot's properties file, by default "./properties.yaml"

### Using Docker

> Soon to come !


# Information

## Contributors

  * [elBichon](https://github.com/elBichon)
  * [Sarlam](https://github.com/sarlam)
  * [Terag](https://github.com/Terag)

## Credits

The bot is developed using the SDK DiscordGo: https://github.com/bwmarrin/discordgo

The bot and characters it quotes are inspired by the universe of Kaamelott: https://kaamelott.com/

## License

Kadok: The code of Kadok is licensed under GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007

Kaamelott: Is a series created by Alexandre Astier, Alain Kappauf and Jean-Yves Robin. All rights reserved to the authors and M6.
