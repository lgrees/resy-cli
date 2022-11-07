# resy-cli <!-- omit in toc -->

<p align="center">
<img src="https://i.ytimg.com/vi/TOecxTy4ZJE/hqdefault.jpg"/>
</p>
<p align="center">
^ a happy <code>resy-cli</code> user getting exactly the reservation he wanted
</p>
<br/>
<br/>

*Disclaimer: This document targets a non-technical audience. For a more technical version of this README, consult [PERUSEME.md]().*

## Table of Contents <!-- omit in toc -->
- [About](#about)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Authentication Setup](#authentication-setup)
- [Scheduling a Booking](#scheduling-a-booking)
- [FAQ](#faq)

## About

`resy-cli` is a program to schedule a resy reservation booking to execute at *exactly* the right time in the future.

You might be asking, why would you do this? This is in fact a great question, as this project is utterly useless in low-demand markets where you're better off booking your reservation through resy's website. However in high-demand markets like NYC, reservation slots at top restaurants are snatched up in a matter of seconds.

After too many instances of losing to people who could click faster than me (or other programs 😅), I decided that enough was enough. While this project initially existed as a simple Node script, I wanted to make something that was easily distributable and usable by friends and family. A few weekends later, `resy-cli` was born.

## Prerequisites

### Terminal Familiarity <!-- omit in toc -->

`resy-cli` is a [command-line interface](https://en.wikipedia.org/wiki/Command-line_interface), and resultingly requires some familiarity with a [terminal emulator]().

If you are using a Mac computer, this [quick introduction]() to using the MacOS default terminal emulator (Terminal) should give you a good foundation to set up `resy-cli` on your computer.

### `at` Permissions <!-- omit in toc -->

Under the hood, `resy-cli` uses another command-line interface called `at` to schedule reservations to book in the future.

While it's unnecessary to have an understanding of how `at` works in order to use `resy-cli`, it _is_ necessary to follow some setup to activate this program on your machine. The following instructions are for MacOS (where `at` comes pre-installed). If you're using a different OS, these steps will vary (and you may have to install `at` separately).

1. Execute the following command from your terminal emulator (this will prompt you for your computer's user credentials):
   ```
   sudo launchctl load -w /System/Library/LaunchDaemons/com.apple.atrun.plist
   ```
2. Navigate to `System Preferences > Security & Privacy > Privacy`
3. Unlock to edit settings
4. Add a new entry to "Full Disk Access"
5. Press `Cmd+Shift+G` at the file selection dialog to select a custom path
6. Add `/usr/libexec/atrun` to the list of commands/applications with "Full Disk Access"

## Installation

To install `resy-cli`, run the following command from your terminal emulator:

```
TODO
```

If the install succeeded, you should see a success message in your terminal.

You can verify that the installation was successful by running the following command:
```
resy -v
```
To view all of the commands that are available to you, run:
```
resy
```

## Authentication Setup
In order to book reservations on your behalf, `resy-cli` needs to know about your resy account credentials. To add these, run:
```
resy setup
```

You'll be prompted to add an api key and an auth token. Follow these steps to find this information:

1. Open a web browser
2. Open Developer Tools (on Chrome: `Chrome > View > Developer > Developer Tools`)
3. Navigate to resy and log in
4. Open the "Network" tab with Developer Tools
5. Search for requests to the domain: `api.resy.com`
6. View the request headers

At this point, you should see something that looks like the following:

The information highlighted in COLOR is what you are looking for.

To verify that this setup was successful, run:

```
resy ping
```
If this command fails, it likely means that the credentials that you provided are incorrect. Repeat the authentication setup again to re-enter your credentials.

Once this command succeeds, you're ready to start booking!

**NOTE:** Resy will periodically expire your credentials. It's a good practice to run `resy ping` before booking to ensure that `resy-cli` can connect to your account.

## Scheduling a Booking

Run:
```
resy schedule
```
Now follow the prompts to schedule your booking.

## FAQ

**Q:** How are my resy credentials stored?<br/>
**A:** Credentials are written to disk on your local machine and are never shared (outside of making requests to resy).

**Q:** What about time zones?<br/>
**A:** `resy-cli` treats the booking time that you input as system times.

As an example, if you want to book at 5:00PM and your computer is set to use the eastern time zone, you will end up attempting to book at 5:00PM eastern.

Reservation times are localized to the location of the restaurant.

**Q:** What if my computer is turned off / asleep at the time when I want to book a reservation?<br/>
**A:** `resy-cli` books reservations using your local machine. If the machine is turned off / asleep, booking will not complete.

**Q:** Are reservation services okay with this?<br/>
**A:** It's up to users to follow the ([Terms of Service](https://resy.com/terms)) set forth by resy. `resy-cli` is just a tool to automate your interactions with resy (no web scraping or anything of the sort).

**Q:** What happens when everyone is using programs like this? Won't this stop being effective?<br/>
**A:** Yes.

**Q:** I think I found a bug, how can I report it?<br/>
**A:** If you have a github account, feel free to submit an issue. If not, feel free to text me.

**Q:** This program _literally_ changed my life, how can I possibly thank the author?<br/>
**A:** Consider taking him out to dinner 🙃.
