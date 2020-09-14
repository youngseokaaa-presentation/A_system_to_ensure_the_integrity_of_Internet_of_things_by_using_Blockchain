// Copyright (c) 2011-2014 The Youngseokcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#ifndef YOUNGSEOKCOIN_QT_YOUNGSEOKCOINADDRESSVALIDATOR_H
#define YOUNGSEOKCOIN_QT_YOUNGSEOKCOINADDRESSVALIDATOR_H

#include <QValidator>

/** Base58 entry widget validator, checks for valid characters and
 * removes some whitespace.
 */
class YoungseokcoinAddressEntryValidator : public QValidator
{
    Q_OBJECT

public:
    explicit YoungseokcoinAddressEntryValidator(QObject *parent);

    State validate(QString &input, int &pos) const;
};

/** Youngseokcoin address widget validator, checks for a valid youngseokcoin address.
 */
class YoungseokcoinAddressCheckValidator : public QValidator
{
    Q_OBJECT

public:
    explicit YoungseokcoinAddressCheckValidator(QObject *parent);

    State validate(QString &input, int &pos) const;
};

#endif // YOUNGSEOKCOIN_QT_YOUNGSEOKCOINADDRESSVALIDATOR_H
