#!/bin/sh

TOPDIR=${TOPDIR:-$(git rev-parse --show-toplevel)}
SRCDIR=${SRCDIR:-$TOPDIR/src}
MANDIR=${MANDIR:-$TOPDIR/doc/man}

YOUNGSEOKCOIND=${YOUNGSEOKCOIND:-$SRCDIR/youngseokcoind}
YOUNGSEOKCOINCLI=${YOUNGSEOKCOINCLI:-$SRCDIR/youngseokcoin-cli}
YOUNGSEOKCOINTX=${YOUNGSEOKCOINTX:-$SRCDIR/youngseokcoin-tx}
YOUNGSEOKCOINQT=${YOUNGSEOKCOINQT:-$SRCDIR/qt/youngseokcoin-qt}

[ ! -x $YOUNGSEOKCOIND ] && echo "$YOUNGSEOKCOIND not found or not executable." && exit 1

# The autodetected version git tag can screw up manpage output a little bit
YSCVER=($($YOUNGSEOKCOINCLI --version | head -n1 | awk -F'[ -]' '{ print $6, $7 }'))

# Create a footer file with copyright content.
# This gets autodetected fine for youngseokcoind if --version-string is not set,
# but has different outcomes for youngseokcoin-qt and youngseokcoin-cli.
echo "[COPYRIGHT]" > footer.h2m
$YOUNGSEOKCOIND --version | sed -n '1!p' >> footer.h2m

for cmd in $YOUNGSEOKCOIND $YOUNGSEOKCOINCLI $YOUNGSEOKCOINTX $YOUNGSEOKCOINQT; do
  cmdname="${cmd##*/}"
  help2man -N --version-string=${YSCVER[0]} --include=footer.h2m -o ${MANDIR}/${cmdname}.1 ${cmd}
  sed -i "s/\\\-${YSCVER[1]}//g" ${MANDIR}/${cmdname}.1
done

rm -f footer.h2m
