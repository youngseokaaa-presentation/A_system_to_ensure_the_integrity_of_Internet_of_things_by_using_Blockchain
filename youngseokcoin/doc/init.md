Sample init scripts and service configuration for youngseokcoind
==========================================================

Sample scripts and configuration files for systemd, Upstart and OpenRC
can be found in the contrib/init folder.

    contrib/init/youngseokcoind.service:    systemd service unit configuration
    contrib/init/youngseokcoind.openrc:     OpenRC compatible SysV style init script
    contrib/init/youngseokcoind.openrcconf: OpenRC conf.d file
    contrib/init/youngseokcoind.conf:       Upstart service configuration file
    contrib/init/youngseokcoind.init:       CentOS compatible SysV style init script

1. Service User
---------------------------------

All three Linux startup configurations assume the existence of a "youngseokcoin" user
and group.  They must be created before attempting to use these scripts.
The OS X configuration assumes youngseokcoind will be set up for the current user.

2. Configuration
---------------------------------

At a bare minimum, youngseokcoind requires that the rpcpassword setting be set
when running as a daemon.  If the configuration file does not exist or this
setting is not set, youngseokcoind will shutdown promptly after startup.

This password does not have to be remembered or typed as it is mostly used
as a fixed token that youngseokcoind and client programs read from the configuration
file, however it is recommended that a strong and secure password be used
as this password is security critical to securing the wallet should the
wallet be enabled.

If youngseokcoind is run with the "-server" flag (set by default), and no rpcpassword is set,
it will use a special cookie file for authentication. The cookie is generated with random
content when the daemon starts, and deleted when it exits. Read access to this file
controls who can access it through RPC.

By default the cookie is stored in the data directory, but it's location can be overridden
with the option '-rpccookiefile'.

This allows for running youngseokcoind without having to do any manual configuration.

`conf`, `pid`, and `wallet` accept relative paths which are interpreted as
relative to the data directory. `wallet` *only* supports relative paths.

For an example configuration file that describes the configuration settings,
see `contrib/debian/examples/youngseokcoin.conf`.

3. Paths
---------------------------------

3a) Linux

All three configurations assume several paths that might need to be adjusted.

Binary:              `/usr/bin/youngseokcoind`  
Configuration file:  `/etc/youngseokcoin/youngseokcoin.conf`  
Data directory:      `/var/lib/youngseokcoind`  
PID file:            `/var/run/youngseokcoind/youngseokcoind.pid` (OpenRC and Upstart) or `/var/lib/youngseokcoind/youngseokcoind.pid` (systemd)  
Lock file:           `/var/lock/subsys/youngseokcoind` (CentOS)  

The configuration file, PID directory (if applicable) and data directory
should all be owned by the youngseokcoin user and group.  It is advised for security
reasons to make the configuration file and data directory only readable by the
youngseokcoin user and group.  Access to youngseokcoin-cli and other youngseokcoind rpc clients
can then be controlled by group membership.

3b) Mac OS X

Binary:              `/usr/local/bin/youngseokcoind`  
Configuration file:  `~/Library/Application Support/Youngseokcoin/youngseokcoin.conf`  
Data directory:      `~/Library/Application Support/Youngseokcoin`
Lock file:           `~/Library/Application Support/Youngseokcoin/.lock`

4. Installing Service Configuration
-----------------------------------

4a) systemd

Installing this .service file consists of just copying it to
/usr/lib/systemd/system directory, followed by the command
`systemctl daemon-reload` in order to update running systemd configuration.

To test, run `systemctl start youngseokcoind` and to enable for system startup run
`systemctl enable youngseokcoind`

4b) OpenRC

Rename youngseokcoind.openrc to youngseokcoind and drop it in /etc/init.d.  Double
check ownership and permissions and make it executable.  Test it with
`/etc/init.d/youngseokcoind start` and configure it to run on startup with
`rc-update add youngseokcoind`

4c) Upstart (for Debian/Ubuntu based distributions)

Drop youngseokcoind.conf in /etc/init.  Test by running `service youngseokcoind start`
it will automatically start on reboot.

NOTE: This script is incompatible with CentOS 5 and Amazon Linux 2014 as they
use old versions of Upstart and do not supply the start-stop-daemon utility.

4d) CentOS

Copy youngseokcoind.init to /etc/init.d/youngseokcoind. Test by running `service youngseokcoind start`.

Using this script, you can adjust the path and flags to the youngseokcoind program by
setting the YOUNGSEOKCOIND and FLAGS environment variables in the file
/etc/sysconfig/youngseokcoind. You can also use the DAEMONOPTS environment variable here.

4e) Mac OS X

Copy org.youngseokcoin.youngseokcoind.plist into ~/Library/LaunchAgents. Load the launch agent by
running `launchctl load ~/Library/LaunchAgents/org.youngseokcoin.youngseokcoind.plist`.

This Launch Agent will cause youngseokcoind to start whenever the user logs in.

NOTE: This approach is intended for those wanting to run youngseokcoind as the current user.
You will need to modify org.youngseokcoin.youngseokcoind.plist if you intend to use it as a
Launch Daemon with a dedicated youngseokcoin user.

5. Auto-respawn
-----------------------------------

Auto respawning is currently only configured for Upstart and systemd.
Reasonable defaults have been chosen but YMMV.
