# GoCampfire

This is a campfire bot in Go. In order for it to work, you'll need to create a user on campfire for it and copy the apikey in My Info when logged in. (To test it, you can just use your own apikey, but it's best to run it as a special user).

In order to use it, you need to create a config file that looks like this:

<pre>
[name]
name = sitename
apikey = sdhfkshghfdgd
room = Room Name Here
</pre>

and save it somewhere as something like campfire.ini

Then to start the bot, run `./campfire ~/location/of/config/file/campfire.ini name` where name is the string you used as the header for the config `[name]`.

## Features - doesn't do any of these yet

- Git stats `gopher buzzard stats`
- Git stats `gopher falcon stats`
- Erlang docs `gopher erlang apply`