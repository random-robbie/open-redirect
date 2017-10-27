# open-redirect
Open Redirect Finder.


[![Capture.png](https://s1.postimg.org/88l48isty7/Capture.png)](https://postimg.org/image/5dsg2qdn6j/)

About
----

Based on the idea from @ak1t4 and his script https://github.com/ak1t4/open-redirect-scanner

This takes 2 files one for the urls to test and one for the payloads.
I've supplied some payloads to try.

It uses the casperjs headless browser to check if the site has redirected to the payload given.
If so it logs the data to found.txt







Status
------

Working for me untested for others



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





To do
--------



Fix any issues reported


append https:// or http:// to urls provided without.
