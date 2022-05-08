# Features

# General

Commands hierarchy:

```
 kadok
 ├─ aide
 ├─ clan
 |  ├─ liste
 |  ├─ rejoindre
 |  └─ quitter
 ├─ groupe
 |  ├─ liste
 |  ├─ rejoindre
 |  └─ quitter
 ├─ aqui
 └─ tatan
```

### Aide

```
kadok aide
```

**Permission:** GetHelp

Returns available actions from the bot. You can also use help with most of the commands to get
more information about its usages.

### Aqui

```
kadok aqui
```

**Permission:** GetCharacterList

Returns available characters for quotes.

### Tatan (Status)

```
kadok tatan
```

**Permission**: _Empty_

Returns general information regarding the bot. Such as build information and licenses.

### Groupe (liste, rejoindre, quitter)

```
kadok groupe <liste|rejoindre|quitter>
```

**Permission**: _Empty_

Used to manage groups of the requesting player. A player can be part of multiple groups.
Only roles with the type `Group` can be used

### Clan (liste, rejoindre, quitter)

```
kadok clan <liste|rejoindre|quitter>
```

**Permission**: _Empty_

Used to manage clans of the requesting player. A player can be part of onlye one clan.
Only roles with the type `Clan` can be used

### Quoting

**Permission:** CallCharacter

Quotes a character if found in one of your messages.
Use `kadok aqui` to get the list of available characters in your guild.
