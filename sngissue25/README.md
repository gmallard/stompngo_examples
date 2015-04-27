# Demo stompngo failures with ackmode client or client-individual #

## Overview ##

The purpose of this subdirectory is to demonstrate failures in the 
current stompngo package at:

* [stompngo at github](https://github.com/gmallard/stompngo)

The demonstrations are intended to:

    1. More fully understand all of the problems with certain ackmodes.
    2. Develop one or more strategies for addressing these problems.
    3. Eventually demonstrate that corrections to stompngo are successful.

In general, these failures currently result in client total blocking (hangs).

For an introduction to this issue, see the discussion details at:

* [stompngo_examples issue 2](https://github.com/gmallard/stompngo_examples/issues/2)

## Basic Test Scenario ##

The test scenarios used here will operate as:

* Place 100 messages on a broker queue
* Subscribe and consume 10 messages
* Attempt normal UNSUBSCRIBE (and/or DISCONNECT)

Details of how to run each scenario are described in comments in individual
source files provided.

