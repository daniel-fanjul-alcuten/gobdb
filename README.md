Description
===========

Database system with the following features:

* Non relational.

That is not ready yet for production because of:

* Non finished.

That will have these others in the future:

* The user defines the model: any go type, it will be held in memory all the time.
* The user defines the operations: any go types, they will read or write the model.
* The user is not expected to access the model but through these operations.
* The write operations are writen to log files using the gob encoding.
* On startup, all operations in the log files will be reapplied in the same order.
* It is responsibility of the user to make the operations deterministic.
* The log files can be compacted in snapshots in the same or different processes.
