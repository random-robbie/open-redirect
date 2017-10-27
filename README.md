# open-redirect
Open Redirect Finder.


Status
------

NOT WORKING AT THE MOMENT.

Based on the idea from @ak1t4 and his script https://github.com/ak1t4/open-redirect-scanner

Requirements
--------------

You need to have nodejs installed.

If you need to install it follow the below:

```
curl -sL https://deb.nodesource.com/setup_8.x | bash -
apt-get install -y nodejs
```


```
apt-get install phantomjs -y
npm -g install casperjs
```


How to run
--------------
```
python redirect urls.txt payloads.txt
```

the urls list must have full http:// or https:// domains and seperated on a new line.

This will test each url with every payload and then write to found.txt when it finds an open redirect.

To do
--------

Fix any issues reported
append https:// or http:// to urls provided without.
