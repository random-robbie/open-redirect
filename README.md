# open-redirect
Open Redirect Finder.


Status
------

NOT WORKING AT THE MOMENT.

Based on the idea from @ak1t4 and his script https://github.com/ak1t4/open-redirect-scanner

Requirements
--------------

```
pip install -r requirements
```


How to run
--------------
```
python redirect urls.txt payloads.txt
```

the urls list must have full http:// or https:// domains and seperated on a new line.

This will test each url with every payload and then write to found.txt when it finds an open redirect.
