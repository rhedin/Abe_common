// Packbge abelogutil is to provide detailed logging while I'm trying to figure out
// what Matthias Ladkau's eliasdb and ecal programs do.
//
// A fair amount of effort is put into logging when the flow arrives in a routine,
// and what the routine is doing.  I don't want to yank it all out at some arbitrary
// moment.  I want to leave it in place, and have it not slow things down in production,
// when the logging is turned off.
//
// Things are arranged so that when the underlog tag is not given at build time,
// the underDisabled.go file is used instead of the underEnabled.go file, and
// the function is empty.  The compiler removes the function from the routine,
// as it is a no-op.
//
// We also supply the UnderEnabled variable.  If there is some more computation
// to be done, in connection with logging, or if the call to UnderPrintf is made
// with arguments that take a lot of time to evaluate, we can put an
// if UnderEnabled { x } around the code, and make it a no-op for elimination
// that way.
//
// The "under" is for understanding.  This is logging for understanding.  More
// detailed than "debug" logging, and even more detailed than (or different than)
// "trace" logging.
package abelogutil
