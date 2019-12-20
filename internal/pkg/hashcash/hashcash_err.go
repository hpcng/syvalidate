// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package hashcash

import "errors"

var HashCashFuturisticErr = errors.New("stamp is futuristic")
var HashCashExpiredErr = errors.New("stamp expired")
var HashCashWrongFormatErr = errors.New("stamp has the wrong format")
var HashCashInsufficientErr = errors.New("stamp does not have enough significant bits equal to zero")
var HashCashSpentErr = errors.New("stamp spent")
