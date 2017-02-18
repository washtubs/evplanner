EV-PLANNER
==========

A solution for programmatically managing your day!

Definitions:
------------

_`happens time`
    Derived from the `event qualifier`_.
_`notification time`
    The time that a reminder notification may be sent. Typically this is inferred, but
    when specified it means "remind me by this time" i.e. "remind me before this time"
_`reminder`
    The thing that gets scheduled. Reminders crucially constist of a description of
    the thing that is supposed to happen, and a time when the thing is supposed to
    happen (see: `happens time`_).
_`event`
    Can be a date or a _`token event`: a parameterized string like
    ``@nearGroceryStore Atlanta``
_`simple date recurrence`
    e.g. Every Monday @1pm, Every first Sunday of the month
_`function date recurrence`
    A function which can be defined taking in the date of the last occurrence, and
    outputting the date of the next.
_`dismissal`
    Dismissal is usually an acknowledgement of the reminder. If supported they can help
    us to not send redundant events and track useful information if the dismissal is
    to mean the user completed the task.

Scheduling Use Cases:
---------------------

* Schedule something to happen once
* Have a simple recurring event, such as once a week. (feature:
  `simple date recurrence`_)
* Have an event which recurs on a changing interval, (e.g. practice jazz
  scales every week, then every month, then every 6 months, etc.) (feature: 
  `function date recurrence`_)
* Events that aren't dates. Events should be able to trigger from things that
  happen. The main hope for this is for the user's phone to send events
  like GPS proximity. (see: `event`_)
* Non-specific times of day. Sometimes I just want something to get done in
  the morning or afternoon and I don't feel like I should have to come up with
  anything specific. Maybe I just want to make sure something gets done in the
  morning sometimes.

Meta Use Cases:
---------------
* _`Projection`, a feature which lets you see what reminders are coming up in the future.
  Also lets you step through the events to see how they play out.
* Be able to do everything with a text editor if you want
* Web and/or phone client
* Tagging support. Might be nice to have an [IMPORTANT] label. I think tags should be
  embedded in the description, so they should be found though description scanning,
  not a special `reminder`_ component.

Not on the radar:
-----------------

Multi-user: Haven't thought about it a lot, so I could be swayed, but I think this sort
of thing makes sense more when you *own* it.

Architecture:
-------------

_`Events server`
    Accepts HTTP requests and runs off a clock. Events server would allow for *arbitrary*
    subscribers. Even shell scripts should be able to run off it.

_`EV-Planner Management server`
    HTTP server which can be used to manage the schedule. Bare bones operations for a
    text editing user: lock (block events processing till unlock), unlock, read, write,
    validate.

_`Notification driven`
    Everything is based on the notification-dismissal workflow, with optional support for
    dismissal (as that can be finicky with different technologies like pushbullet).
    
**Flat file storage**: the schedule will be storeable as a flat file for cold 
storage. [1]_


.. [1] Most of the time the server will be running and just have the whole thing loaded
   as an object. If we rely on flat file storage just for persistence across shutdowns,
   there won't be a need for performance there, and we can reuse the code for formatting
   and parsing.

Syntax
------

Parser should have very permissive date parsing

Sample 1::

    Jul 04 2017
        @nearGroceryStore
        -- Pick up hot dogs!
        @1:00pm
        -- Get fireworks
        @getHome -- Blow up fireworks

There are three reminders here: pick up the hot dogs, get the fireworks, and blow up the
fireworks. Note that only one of them is a time, the other two are examples of a
`token event`_. The meaning is as follows. Any time during July 4th, if I go near the
grocery store, I will get a reminder saying to pick up the hot dogs. Sometimes before 1:00
I'll get a reminder to get the fireworks. And if I arrive home at any time after 1 on
July 4th, I'll get a reminder to blow them up.

Sample 2::
    
    Mon. Sept 6 2017
        ...
    @nearGroceryStrore -- Pick up milk (recurrence-spec-blah)
    Tues. Sept 7 2017
    Wed. Sep 8 2017
        @1:00pm -- Eat donuts
    @nearGroceryStore -- Pick up milk
    Tues. Sept 14 2017
        @1:00pm -- Eat pizza

Generally if a `token event`_ reminder is to recur, it should be placed in `rank-0`_. Any
higher and it is possible the reminder will never fire, being lost because the day went
by during which it was supposed to occur.

.. _`rank-0`:
.. _`rank-1`:

_`event qualifier`
    An event in sum represents exactly one point in time, but it can be broken up into
    components, known as qualifiers. Such as the date and time, or perhaps even the
    month, year, day and time. For simplicity's sake I would like to keep it as just
    the two. The date is a *rank-0* qualifier, and the time is a *rank-1* qualifier. A
    `token event`_ can be either *rank-0* or *rank-0*

