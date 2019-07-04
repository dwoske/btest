
This is a test program to help understand 
issue #908 with Badger using 100% of CPU.  This
happens after a process using a managed badger
database stops without closing the DB, then
that process is restarted.

https://github.com/dgraph-io/badger/issues/908

To test this, start the process and watch the CPU
usage.

Hit the localhost:9999/stop endpoint.
Start the process again, and watch cpu usage.

Continue.

Upon about the 3rd restart of the process, the
CPU pegs at over 100% and does not drop back to
low levels.


The profile for the problem looks like this:

```
      flat  flat%   sum%        cum   cum%
    15.62s 46.57% 46.57%     33.52s 99.94%  github.com/dgraph-io/badger/y.(*WaterMark).process.func1
     9.04s 26.95% 73.52%      9.72s 28.98%  runtime.mapaccess1_fast64
     7.15s 21.32% 94.84%      8.18s 24.39%  runtime.mapdelete_fast64
     1.71s  5.10% 99.94%      1.71s  5.10%  runtime.newstack
         0     0% 99.94%     33.52s 99.94%  github.com/dgraph-io/badger/y.(*WaterMark).process
```

Will update this once the issue is resolved.


