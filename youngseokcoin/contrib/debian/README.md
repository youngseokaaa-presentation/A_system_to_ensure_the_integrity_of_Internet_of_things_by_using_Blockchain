
Debian
====================
This directory contains files used to package youngseokcoind/youngseokcoin-qt
for Debian-based Linux systems. If you compile youngseokcoind/youngseokcoin-qt yourself, there are some useful files here.

## youngseokcoin: URI support ##


youngseokcoin-qt.desktop  (Gnome / Open Desktop)
To install:

	sudo desktop-file-install youngseokcoin-qt.desktop
	sudo update-desktop-database

If you build yourself, you will either need to modify the paths in
the .desktop file or copy or symlink your youngseokcoin-qt binary to `/usr/bin`
and the `../../share/pixmaps/youngseokcoin128.png` to `/usr/share/pixmaps`

youngseokcoin-qt.protocol (KDE)

