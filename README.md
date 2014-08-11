# buildbox-lifx

This is a small service to monitor buildbox builds and change the lifx when something fails.

# Bulding

The following command will compile an executable and put it in `bin/buildbox-lifx`.

```
make
```

If you wanna run this on a beagle bone or the like device.

```
GOARCH=arm GOOS=linux make
```

Or for the raspberry pi.

```
GOARM=5 GOARCH=arm GOOS=linux make
```


# Usage

```
NAME:
   buildbox-lifx - Monitors buildbox and changes lifx bulbs to reflect success or failure

USAGE:
   buildbox-lifx [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --apikey 		buildbox api key
   --branch 'master'	branch to filter builds
   --bulb 'build'	the label of the bulb you want to control
   --debug		enable debug logging
   --help, -h		show help
   --version, -v	print the version
```

So to get started.

1. Setup your globe on the same wifi as your monitor box.
2. Using the lifx app label a globe `build`
3. Grab your person api key from http://buildbox.io
4. Start the service passing in your key.

```
buildbox-lifx -apikey=XXXXX
```

# Disclaimer

This is currently very early release, everything can and will change.

# License

Copyright (c) 2014 Mark Wolfe
Licensed under the MIT license.