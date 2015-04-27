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

For an introduction to this issue, see the initial discussion details at:

* [stompngo_examples issue 2](https://github.com/gmallard/stompngo_examples/issues/2)

The primary issue is filed at:

* [stompngo issue 25](https://github.com/gmallard/stompngo/issues/25)

## Basic Test Scenario ##

The test scenarios used here will operate as:

* Place 100 messages on a broker queue
* Subscribe and consume 10 messages
* Attempt normal UNSUBSCRIBE (and/or just DISCONNECT)

Details of how to run each scenario are described in comments in individual
source files provided.

## Current Failure Examples ##

See the run.sh file in each subdirectory for a detailed method of
reproducing the problem.

Each subdirectory also contains stack traces in each hang (or in several
cases a pure panic).  Stack traces are provided for brokers ActiveMQ,
Apollo, and RabbitMQ.  Note that in some cases, these brokers yield differing
results.

These failure recreations were run with stompngo at commit 048c068.

Failure examples in each subdirectory are as follows:

* client01-hang - SUBSCRIBE with ack mode client, attempt UNSUBSCRIBE
* client02-hang - SUBSCRIBE with ack mode client, attempt immediate DISCONNECT (no UNSUBSCRIBE)
* clientindividual01-hang - SUBSCRIBE with ack mode client-individual, attempt UNSUBSCRIBE
* clientindividual02-hang - SUBSCRIBE with ack mode client-individual, attempt immediate DISCONNECT (no UNSUBSCRIBE)

