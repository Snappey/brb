## BRB (Be Right Back) Discord Bot

Time how long your "friends" actually are when they say brb. (it's basically just a timer)

### Usage

Two commands are exposed `brb` and `back`, these are invoked using Slash commands or mentioning the bot.

#### brb

Marks the author as brb with a default expected time of `5 minutes`. 

```text
/brb duration: 15m // Mark yourself as brb with a target duration of 15 minutes

/brb user: @friend // Mark your friend as brb with a default of 5 minutes
```

#### back

Marks the author as back if they have an active brb session.

```text
/back // Will mark yourself as back if you have an active brb session

/back user: @friend // Will mark your friend as back if they have an active brb session
```

#### Mentions

Mentioning `@brb` will toggle, if a user has an active brb session they will be set as back and vice versa.

Optionally you can target other users via mentioning them and change the expected time by providing a formatted time string.

```text
@brb 15m // Will mark yourself as brb with a target of 15 minutes

@brb @friend // Will mark your friend as brb with the default of 5 minutes
```

### Setup

Configured via Environment Variables, you can get your AppId / AuthToken from [here](https://discord.com/developers/applications) by creating a new application.

```html
APP_ID     = <discord_bot_app_id>
AUTH_TOKEN = <discord_bot_secret>
    
OPTIONAL:
GUILD_ID = <testing_guild_id>
```

Providing a GuildId will create guild Slash commands instead of global.

#### Docker

You can get images from GitHub or pull the latest from `docker pull ghcr.io/snappey/brb:latest`

##### Example

```shell
docker run -it \
  -e APP_ID=discord_bot_app_id \
  -e AUTH_TOKEN=discord_bot_secret \
   ghcr.io/snappey/brb:latest
```

or

```shell
docker run -it \
  -e APP_ID=discord_bot_app_id \
  -e AUTH_TOKEN=discord_bot_secret \
  -e GUILD_ID=testing_guild_id \
   ghcr.io/snappey/brb:latest
```
