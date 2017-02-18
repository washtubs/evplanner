HTTP API for editing schedule files
===================================

We can lock for editing, preventing the scheduler from firing an event and modifying the
state while we are editing.

The basic algorithm for editing should be as follows::
    
    Acquire the write lock to prevent editing by automated processes.
    Read to a temp file
    ... edit the file ...
    Write from the temp file
    Release the write lock

If you take too long to edit the file, and the write lock releases automatically, it's 
possible that a process has modified the schedule file, but it's also possible it hasn't.
This should be handled as follows::

    Acquire the write lock to prevent editing by automated processes.
    Read to a temp file
    ... edit the file but take a while ...
    Write from the temp file
        Received an error indicating the lock is not acquired
    Read the file to another temp file.
    If the contents differ:
        Allow the user to merge them to a new temp file
        Acquire the write lock again
        Write from the merged temp file
    Else
        Acquire the write lock again
        Write from the original temp file
    Release the write lock


.. _`lock`:

``PUT /lock``
==============

Obtains a write lock on the schedule file so that any operation that would need to edit
the state of the file will need to wait until this lock is released. Should be invoked
with a timeout so as not to accidentally leave unlocked indefinitely. A max timeout 
should also be defined, (30 minutes).

.. _`unlock`:

``PUT /unlock``
===============

Releases the write lock notifying the scheduler if it was waiting, that it can start 
firing events again.

.. _`write`:

``PUT /write``
==============

Writes the schedule file. Expects the lock to have already been obtained or sends an 
error indicating it is released.

.. _`read`:

``GET /read``
==============

Reads the schedule file. Does not require any lock.
