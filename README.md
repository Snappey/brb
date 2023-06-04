## BRB (Be Right Back) Discord Bot

Time how long your "friends" actually are when they say brb. (it's basically just a timer)

### Usage

Three commands are exposed `brb`, `back` and `gonefor`, these are invoked using Slash commands or mentioning the bot.

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

When a user has returned buttons will be created asking for another user to confirm, this lasts for 5 minutes before timing out if no one confirms.

- If a user was marked as brb by another person, they are the only user who can confirm for the first 2 minutes.
  - If they are confirmed by another user the reporting user is marked as brb.
- If a user's confirmation timesout the brb is finished as if it was confirmed.
- If a user's confirmation is rejected the brb carries on.
- While a brb is awaiting confirmation it does not increase the running duration.

#### Mentions

Mentioning `@brb` will toggle, if a user has an active brb session they will be set as back and vice versa.

Optionally you can target other users via mentioning them and change the expected time by providing a formatted time string.

```text
@brb 15m // Will mark yourself as brb with a target of 15 minutes

@brb @friend // Will mark your friend as brb with the default of 5 minutes
```

#### Gone For

Gets the ongoing duration of a users brb.

```text
/gonefor @friend
```

### Metrics

Prometheus metrics are exported at `:8080/metrics`.

```yaml

brb_session_active: 
    help: currently active sessions
    labels: guild_id, user_id
    type: gauge

brbSessionDuration:
    help: brb duration from start to finish recorded in seconds
    labels: guild_id, user_id
    type: histogram

brbSessionTargetDuration:
    help: brb target duration set by the reporting user recorded in seconds
    labels: guild_id, user_id
    type: histogram

brbSessionLateDifferenceDuration:
    help: brb difference between finished duration and target duration >0 recorded in seconds
    labels: guild_id, user_id
    type: histogram

brbSessionEarlyDifferenceDuration:
    help: brb difference between finished duration and target duration <0 recorded in seconds
    labels: guild_id, user_id
    type: histogram

bucketDistribution:
  1 minute
  5 minutes
  15 minutes
  30 minutes
  1 hour
  3 hours
  6 hours
  12 hours
  1 day
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
