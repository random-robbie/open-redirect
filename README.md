# open-redirect
Open Redirect Finder.


[![Capture.png](https://s1.postimg.org/88l48isty7/Capture.png)](https://postimg.org/image/5dsg2qdn6j/)

Status
------

Working for me untested for others

Based on the idea from @ak1t4 and his script https://github.com/ak1t4/open-redirect-scanner

Requirements
--------------

You need to have nodejs installed.

If you need to install it follow the below:

```
curl -sL https://deb.nodesource.com/setup_8.x | bash -
apt-get install -y nodejs
```



Normal Requirements

```
apt-get install phantomjs xvfb -y
npm -g install casperjs
```


How to run
--------------


```
python redirect.py urls.txt payloads.txt
```

the urls list must have full http:// or https:// domains and seperated on a new line.

This will test each url with every payload and then write to found.txt when it finds an open redirect.

To do
--------



Fix any issues reported


append https:// or http:// to urls provided without.
