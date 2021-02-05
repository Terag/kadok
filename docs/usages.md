# Features

# General

Commands hierarchy:

```
 kadok
 ├─ aqui
 ├─ ping
 ├─ pong
 └─ tatan
```

### Help

```
kadok help
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

### Ping / Pong

```
kadok <ping|pong>
```

**Permission**: _Empty_

Returns ping or pong.

### Quoting

**Permission:** CallCharacter

Quotes a character if found in one of your messages.
Use `kadok aqui` to get the list of available characters in your guild.