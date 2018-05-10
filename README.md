# Why Awake?
An OSX menu bar app to let you know when an application is keeping your computer from sleeping, and force it to not sleep for a period.

## Install

* [Download the app](https://github.com/caseymrm/whyawake/releases/download/v0.2/WhyAwake.app.zip)
* Unzip it
* Put it in Applications
* Run it
* To run every time you login, select "start at login" from the menu bar

## Screenshots

### Your Mac can sleep normally

![Can Sleep](https://github.com/caseymrm/whyawake/raw/master/static/cansleep.png)

### Your Mac won't sleep because a webpage is playing a video

![Can't Sleep](https://github.com/caseymrm/whyawake/raw/master/static/cantsleep.png)

### You have chosen to keep your Mac from sleeping for 10 minutes

![Prevented](https://github.com/caseymrm/whyawake/raw/master/static/prevented.png)

## Built with

* [Menuet](https://github.com/caseymrm/menuet) - build menu bar apps in Go
* [go-assertions](https://github.com/caseymrm/go-assertions) - detect when your Mac can sleep
* [go-caffeinate](https://github.com/caseymrm/go-caffeinate) - keep your Mac from sleeping

## License

MIT

