// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build BOOST_STACKTRACE_USE_BACKTRACE

package cgosymbolizer

/*
#cgo windows CXXFLAGS:
#cgo !windows CXXFLAGS: -DBOOST_STACKTRACE_USE_BACKTRACE
#cgo !windows LDFLAGS: -ldl -lbacktrace
*/
import "C"

// https://www.boost.org/doc/libs/develop/doc/html/stacktrace/configuration_and_build.html
